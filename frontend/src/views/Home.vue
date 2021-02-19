<template>
  <div class="h-screen w-screen flex flex-col">
    <header class="p-2 flex justify-end shadow-sm">
      <el-button type="primary" round>Use Images</el-button>
      <el-button type="primary" icon="el-icon-s-tools" circle></el-button>
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
          class="p-2 w-56 h-56 bg-white rounded-xl shadow-sm flex flex-col cursor-pointer ml-2 mt-2 hover:shadow-md"
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

export default defineComponent({
  name: "Home",
  components: {
    DragFile
  },
  setup() {
    const dragShow = ref(true);
    const files = ref<File[]>([]);
    const filesData = ref<FileData[]>([]);
    const dragChange = (fs: FileList) => {
      files.value = [...files.value, ...[].slice.apply(fs)];
    };
    watch(files, async () => {
      const timeStart = new Date().getTime();
      const f: FileData[] = [];
      for (const v of files.value) {
        const src = await readAsDataURL(v);
        await window.backend.Manager.HandleFile(
          JSON.stringify({
            name: v.name,
            size: v.size,
            data: src.split(",")[1]
          })
        );
        f.push({ name: v.name, size: v.size, src });
        const timeEnd = new Date().getTime();

        console.log("timeEnd - timeStart: %s", timeEnd - timeStart);
      }
      filesData.value = [...filesData.value, ...f];
    });
    return { filesData, dragShow, dragChange };
  }
});
</script>

<style scoped lang="stylus"></style>
