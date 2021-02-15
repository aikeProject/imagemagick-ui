<template>
  <div class="h-screen w-screen flex flex-col">
    <drag-file class="flex-grow" @change="dragChange"></drag-file>
    <div class="m-2">
      <el-table
        :class="{ hidden: !filesData.length }"
        :data="filesData"
        style="width: 100%"
        max-height="500"
      >
        <el-table-column prop="name" label="文件名" width="180">
        </el-table-column>
        <el-table-column prop="size" label="文件大小"> </el-table-column>
      </el-table>
    </div>
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
    const files = ref<File[]>([]);
    const dragChange = fs => {
      files.value = [...files.value, ...[].slice.apply(fs)];
    };
    const filesData = computed(() => {
      return files.value.map(v => ({ name: v.name, size: v.size }));
    });
    return { filesData, dragChange };
  }
});
</script>

<style scoped lang="stylus"></style>
