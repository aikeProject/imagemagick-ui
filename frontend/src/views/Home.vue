<template>
  <div class="h-screen w-screen flex">
    <drag-file v-model:show="drag" @change="dragChange"></drag-file>
  </div>
</template>

<script>
import { defineComponent, ref } from "vue";
import DragFile from "components/DragFile";

export default defineComponent({
  name: "Home",
  components: {
    DragFile
  },
  setup() {
    const drag = ref(true);
    const handleBasic = async () => {
      const s = await window.backend.basic();
      alert(s);
    };
    const HandleResize = async () => {
      await window.backend.HandleResize();
    };
    const handleFileChange = (file, fileList) => {
      console.log(file);
      console.log(fileList);
    };
    const dragChange = files => {
      console.log(files);
    };

    return { drag, handleBasic, HandleResize, handleFileChange, dragChange };
  }
});
</script>

<style scoped lang="stylus"></style>
