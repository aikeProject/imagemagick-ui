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
        class="flex flex-wrap flex-grow self-start p-4"
      >
        <div
          class="p-2 h-36 bg-white rounded-xl shadow-sm flex flex-grow cursor-pointer mb-3 border border-gray-300"
          v-for="item in filesData"
          :key="item.name"
        >
          <el-image class="rounded w-36" :src="item.src" fit="cover"></el-image>
          <div class="relative space-y-1 pl-4 flex-grow text-gray-500">
            <div class="text-sm truncate">文件名: {{ item.name }}</div>
            <div class="text-sm">大小: {{ item.size }}</div>
            <div class="text-sm">类型: jpeg</div>
            <div class="absolute inset-x-0 bottom-0 pl-4">
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
import DragFile from "components/DragFile.vue";
import { readAsDataURL } from "lib/filw";
import { FileData } from "views/Home";
import { FileStatus } from "common/enum";
import { EventFps } from "lib/event-fps";
export default defineComponent({
  name: "Home",
  components: {
    DragFile
  },
  setup() {
    const dragShow = ref(true);
    const filesData = ref<FileData[]>([]);
    const fileSpeed = ref(100 * 1000);

    watch(filesData, function(v) {
      console.log("filesData update");
      // 清空操作之后，显示拖拽区域
      if (!v.length) dragShow.value = true;
    });

    onMounted(() => {
      let time = 0;
      let last = 0;
      const fileTimeMap: { [key: string]: number } = {};
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
            fileTimeMap[v.id] = (fileTimeMap[v.id] || 0) + f;
            if (v.status === FileStatus.Start) {
              v.progress = parseFloat(
                Math.min(
                  ((fileSpeed.value * fileTimeMap[v.id]) / v.size) * 100,
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

    // 拖拽选择文件
    const dragChange = async (fs: FileList) => {
      const files: File[] = [].slice.apply(fs);
      for (const v of files) {
        const timeStart = new Date().getTime();
        const src = await readAsDataURL(v);
        const f: FileData = {
          id: createFileId(v.name, v.size),
          name: v.name,
          size: v.size,
          src,
          status: FileStatus.NotStarted,
          progress: 0
        };
        filesData.value.push(f);
        const timeEnd = new Date().getTime();
        console.log("base64 %s => %d ms", v.name, timeEnd - timeStart);
      }

      for (const v of filesData.value) {
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
        const speed = (v.size / (timeEnd - timeStart)) * 1000;
        fileSpeed.value = speed;
        console.log(
          "file %s \ntime %d ms \nsize => %d byte \nspeed => %fkb/s",
          v.name,
          timeEnd - timeStart,
          v.size,
          speed / 1024
        );
      }
    };

    const handleConvert = async () => {
      const { Convert } = window.backend.Manager;
      await Convert();
    };

    const handleClear = async () => {
      const { Clear } = window.backend.Manager;
      filesData.value = [];
      await Clear();
    };

    return { filesData, dragShow, dragChange, handleConvert, handleClear };
  }
});
</script>

<style scoped lang="stylus"></style>
