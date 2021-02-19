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
          class="p-2 w-56 h-56 bg-white rounded-xl shadow-sm flex flex-col cursor-pointer ml-2 mt-2 border border-gray-300"
          v-for="item in filesData"
          :key="item.name"
        >
          <el-image class="rounded" :src="item.src" fit="cover"></el-image>
          <div class="space-y-1 pt-2 text-gray-500">
            <div class="text-sm truncate">文件名: {{ item.name }}</div>
            <div class="text-sm">大小: {{ item.size }}</div>
            <div class="text-sm">类型: jpeg</div>
          </div>
        </div>
      </figure>
    </main>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, watch } from "vue";
import DragFile from "components/DragFile.vue";
import { readAsDataURL } from "lib/filw";
import { FileData } from "views/Home";
import { FileStatus } from "common/enum";

export default defineComponent({
  name: "Home",
  components: {
    DragFile
  },
  setup() {
    const dragShow = ref(true);
    const filesData = ref<FileData[]>([]);

    watch(filesData, function(v) {
      // 清空操作之后，显示拖拽区域
      if (!v.length) dragShow.value = true;
    });

    // 拖拽选择文件
    const dragChange = async (fs: FileList) => {
      const files: File[] = [].slice.apply(fs);
      for (const v of files) {
        const timeStart = new Date().getTime();
        const src = await readAsDataURL(v);
        const f: FileData = {
          name: v.name,
          size: v.size,
          src,
          status: FileStatus.NotStarted
        };
        filesData.value.push(f);
        const timeEnd = new Date().getTime();
        console.log("base64 %s => %d ms", v.name, timeEnd - timeStart);
      }
      for (const v of filesData.value) {
        const timeStart = new Date().getTime();
        await window.backend.Manager.HandleFile(
          JSON.stringify({
            name: v.name,
            size: v.size,
            data: v.src.split(",")[1],
            status: FileStatus.NotStarted
          })
        );
        const timeEnd = new Date().getTime();
        console.log("file %s => %d ms", v.name, timeEnd - timeStart);
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
