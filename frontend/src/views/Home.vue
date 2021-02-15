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
          <el-image
            class="rounded"
            src="https://fuss10.elemecdn.com/e/5d/4a731a90594a4af544c0c25941171jpeg.jpeg"
            fit="cover"
          ></el-image>
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
import { defineComponent, ref, computed } from "vue";
import DragFile from "components/DragFile";

export default defineComponent({
  name: "Home",
  components: {
    DragFile
  },
  setup() {
    const dragShow = ref(true);
    const files = ref<File[]>([]);
    const dragChange = fs => {
      files.value = [...files.value, ...[].slice.apply(fs)];
    };
    const filesData = computed(() => {
      return files.value.map(v => ({ name: v.name, size: v.size }));
    });
    return { filesData, dragShow, dragChange };
  }
});
</script>

<style scoped lang="stylus"></style>
