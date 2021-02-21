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
        v-show="!!filesView.length"
        class="flex flex-wrap flex-col flex-grow self-start p-4"
      >
        <div
          class="p-2 h-28 bg-white rounded-md shadow-sm flex flex-grow cursor-pointer mb-3 border border-gray-200"
          v-for="item in filesView"
          :key="item.name"
        >
          <el-image class="rounded w-28" :src="item.src" fit="cover"></el-image>
          <div class="relative space-y-1 pl-3 flex-grow text-gray-500">
            <div class="text-sm truncate">文件名: {{ item.name }}</div>
            <div class="text-sm">大小: {{ item.size }}</div>
            <div class="text-sm">类型: jpeg</div>
            <div v-show="item.show" class="absolute inset-x-0 bottom-0 pl-3">
              <el-progress
                :percentage="item.progress"
                :status="item.statusStr"
              ></el-progress>
            </div>
          </div>
        </div>
      </figure>
    </main>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, watch, computed, onMounted } from "vue";
import { ElMessage } from "element-plus";
import Wails from "@wailsapp/runtime";
import DragFile from "components/DragFile.vue";
import { readAsDataURL } from "lib/filw";
import { Complete, FileData } from "views/Home";
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

    // 根据id查询文件实例
    const getFileById = (id: string) => {
      if (!filesData.value.length) return;
      return filesData.value.find(v => v.id === id);
    };

    /**
     * 根据文件名和大小创建文件id
     * @param name
     * @param size
     */
    const createFileId = (name: string, size: number) => {
      return name + size.toString();
    };

    // 向golang程序发送文件数据
    const sendFile = async (files: FileData[]) => {
      files = files.filter(v => v.status === FileStatus.NotStarted);
      for (const v of files) {
        // 数据已发送至golang处理程序 无需重复发送
        if (v.status === FileStatus.SendSuccess) continue;
        const timeStart = new Date().getTime();
        // 开始发送
        v.status = FileStatus.Start;
        await window.backend.Manager.HandleFile(
          JSON.stringify({
            id: v.id,
            name: v.name,
            size: v.size,
            data: v.src.split(",")[1],
            status: FileStatus.Start
          })
        );
        // 发送完成
        v.status = FileStatus.SendSuccess;
        v.progress = 100;
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
    const covertFile = async (v: File) => {
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
    const covertFileData = async (fs: FileList): Promise<FileData[]> => {
      const files: File[] = [].slice.apply(fs);
      if (!files.length) return [];
      return (await Promise.all(files.map(v => covertFile(v)))).filter(
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
      const files = await covertFileData(fs);
      filesData.value = [...filesData.value, ...files];
      await sendFile(filesData.value);
    };

    // 调用golang程序处理文件
    const handleConvert = async () => {
      // 改变文件状态
      filesData.value.forEach(v => {
        if (v.status == FileStatus.SendSuccess) {
          v.status = FileStatus.Running;
          v.progress = 0;
        }
      });
      fileTimeMap.value = {};
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

    const filesView = computed(() => {
      return filesData.value.map(v => ({
        ...v,
        // show 进度条是否显示
        show: v.status === 4 || (v.progress < 100 && v.progress > 0)
      }));
    });

    onMounted(() => {
      let time = 0;
      let last = 0;
      // fps
      EventFps.on<number>("update", function(f) {
        if (filesData.value.every(v => v.progress >= 100)) return;
        // 所有文件都已经传输至golang
        if (filesData.value.every(v => v.status === FileStatus.SendSuccess)) {
          return;
        }
        // 所有文件都经由golang处理完毕
        if (filesData.value.every(v => v.status === FileStatus.Done)) {
          return;
        }
        f = f || 0;
        time += f;
        if ((time - last) * 1000 > 100) {
          for (const v of filesData.value) {
            if (
              v.status === FileStatus.Start ||
              v.status === FileStatus.Running
            ) {
              fileTimeMap.value[v.id] = (fileTimeMap.value[v.id] || 0) + f;
              v.progress = parseFloat(
                Math.min(
                  ((fileSpeed.value * fileTimeMap.value[v.id]) / v.size) * 100,
                  99
                ).toFixed(1)
              );
            }
          }
          last = time;
        }
      });

      // file:complete events 文件处理完成后收到的数据
      Wails.Events.On("file:complete", (data: Complete) => {
        if (!data) return;
        const file = getFileById(data.id);
        if (file) {
          // 更新文件状态
          file.status = data.status;
          file.statusStr = "success";
          file.progress = 100;
        }
      });
    });

    return {
      filesView,
      dragShow,
      dragChange,
      handleConvert,
      handleClear
    };
  }
});
</script>

<style scoped lang="stylus"></style>
