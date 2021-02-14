<template>
  <div class="h-screen w-screen flex">
    <drag-file></drag-file>
  </div>
</template>

<script>
import { defineComponent } from "vue";
import DragFile from "components/DragFile";

export default defineComponent({
  name: "Home",
  components: {
    DragFile
  },
  setup() {
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

    return { handleBasic, HandleResize, handleFileChange };
  }
});
</script>

<style scoped lang="stylus"></style>
