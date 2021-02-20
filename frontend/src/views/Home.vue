<template>
  <div class="h-screen w-screen flex flex-col">
    <header class="p-2 flex justify-end">
      <el-button @click="handleConvert" type="primary" round>Convert</el-button>
      <el-button @click="handleClear" type="primary" round>清空</el-button>
      <el-button type="primary" round>设置</el-button>
    </header>
    <main class="flex-1 flex">
      <drag-file
        v-model:show="dragShow"
        class="flex-grow"
        @change="dragChange"
      ></drag-file>
      <figure
        v-show="!!filesData.length"
        class="flex flex-wrap flex-col flex-grow self-start p-4"
      >
        <div
          class="p-2 h-28 bg-white rounded-md shadow-sm flex flex-grow cursor-pointer mb-3 border border-gray-200"
          v-for="item in filesData"
          :key="item.name"
        >
          <el-image class="rounded w-28" :src="item.src" fit="cover"></el-image>
          <div class="relative space-y-1 pl-3 flex-grow text-gray-500">
            <div class="text-sm truncate">文件名: {{ item.name }}</div>
            <div class="text-sm">大小: {{ item.size }}</div>
            <div class="text-sm">类型: jpeg</div>
            <div
              v-show="item.progress < 100 && item.progress > 0"
              class="absolute inset-x-0 bottom-0 pl-3"
            >
              <el-progress :percentage="item.progress"></el-progress>
            </div>
          </div>
        </div>
      </figure>
    </main>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, watch, onMounted } from "vue";
import { ElMessage } from "element-plus";
import DragFile from "components/DragFile.vue";
import { readAsDataURL } from "lib/filw";
import { FileData } from "views/Home";
import { FileStatus } from "common/enum";
import { EventFps } from "lib/event-fps";
import useDrag from "composables/useDrag";
export default defineComponent({
  name: "Home",
  components: {
    DragFile
  },
  setup() {
    const dragRef = ref(document.body);
    const dragShow = ref(true);
    const filesData = ref<FileData[]>([]);
    const fileSpeed = ref(1000 * 1000);
    const fileTimeMap = ref<{ [key: string]: number }>({});
    const { files } = useDrag(dragRef);

    onMounted(() => {
      let time = 0;
      let last = 0;
      // fps
      EventFps.on<number>("update", function(f) {
        if (filesData.value.every(v => v.progress >= 100)) return;
        // 所有文件都已经传输至golang
        if (filesData.value.every(v => v.status === FileStatus.SendSuccess)) {
          filesData.value.forEach(v => (v.progress = 100));
          return;
        }
        // 所有文件都经由golang处理完毕
        if (filesData.value.every(v => v.status === FileStatus.Done)) {
          filesData.value.forEach(v => (v.progress = 100));
          return;
        }
        f = f || 0;
        time += f;
        if ((time - last) * 1000 > 100) {
          for (const v of filesData.value) {
            if (v.status === FileStatus.Start) {
              fileTimeMap.value[v.id] = (fileTimeMap.value[v.id] || 0) + f;
              v.progress = parseFloat(
                Math.min(
                  ((fileSpeed.value * fileTimeMap.value[v.id]) / v.size) * 100,
                  99
                ).toFixed(1)
              );
            } else if (v.status === FileStatus.SendSuccess) {
              v.progress = 100;
            }
          }
          last = time;
        }
      });
    });

    /**
     * 根据文件名和大小创建文件id
     * @param name
     * @param size
     */
    const createFileId = (name: string, size: number) => {
      return name + size.toString();
    };

    // 向golang程序发送文件数据
    const SendFile = async (files: FileData[]) => {
      for (const v of files) {
        // 数据已发送至golang处理程序 无需重复发送
        if (v.status === FileStatus.SendSuccess) continue;
        const timeStart = new Date().getTime();
        // 开始发送
        v.status = FileStatus.Start;
        await window.backend.Manager.HandleFile(
          JSON.stringify({
            name: v.name,
            size: v.size,
            data: v.src.split(",")[1],
            status: FileStatus.Start
          })
        );
        // 发送完成
        v.status = FileStatus.SendSuccess;
        const timeEnd = new Date().getTime();
        fileSpeed.value = (v.size / (timeEnd - timeStart)) * 1000;
        console.log(
          "file %s \ntime %d ms \nsize => %d byte \nspeed => %fkb/s",
          v.name,
          timeEnd - timeStart,
          v.size,
          Number((fileSpeed.value / 1024).toFixed(2))
        );
      }
    };

    // 将文件装换位base64字符串
    const CovertFile = async (v: File) => {
      const file: FileData = {
        id: createFileId(v.name, v.size),
        name: v.name,
        size: v.size,
        src: "",
        status: FileStatus.NotStarted,
        progress: 0
      };
      try {
        file.src = await readAsDataURL(v);
        return file;
      } catch (e) {
        ElMessage({
          message: `名为"${v.name}"的文件转换base64失败`,
          type: "error"
        });
      }
      return file;
    };

    // 构造FileData数据
    const CovertFileData = async (fs: FileList): Promise<FileData[]> => {
      const files: File[] = [].slice.apply(fs);
      if (!files.length) return [];
      return (await Promise.all(files.map(v => CovertFile(v)))).filter(
        v => v.src
      );
    };

    // 检查是否有文件正在发送中
    const checkSend = () =>
      filesData.value.some(
        v => v.status === FileStatus.NotStarted || v.status == FileStatus.Start
      );

    // 拖拽选择文件，并将文件发送至golang程序
    const dragChange = async (fs: FileList) => {
      if (!fs) return;
      const files = await CovertFileData(fs);
      filesData.value = [...filesData.value, ...files];
      await SendFile(filesData.value);
    };

    // 调用golang程序处理文件
    const handleConvert = async () => {
      const { Convert } = window.backend.Manager;
      await Convert();
    };

    // 清空
    const handleClear = async () => {
      const { Clear } = window.backend.Manager;
      filesData.value = [];
      fileTimeMap.value = {};
      await Clear();
    };

    watch(filesData, function(v) {
      // 清空操作之后，显示拖拽区域
      if (!v.length) dragShow.value = true;
    });

    watch(files, function() {
      if (checkSend()) {
        ElMessage({
          message: "等待中...",
          type: "warning"
        });
        return;
      }
      // 继续添加文件
      files.value && dragChange(files.value);
    });

    return {
      filesData,
      dragShow,
      dragChange,
      handleConvert,
      handleClear
    };
  }
});
</script>

<style scoped lang="stylus"></style>
