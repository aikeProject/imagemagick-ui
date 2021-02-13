<template>
  <div class="p-8">
    <el-button @click="handleBasic" type="primary">Basic</el-button>
    <el-button @click="HandleResize" type="primary">HandleResize</el-button>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  name: "Home",
  components: {},
  setup() {
    const handleBasic = async () => {
      const s = await window.backend.basic();
      alert(s);
    };
    const HandleResize = async () => {
      await window.backend.HandleResize();
    };

    return { handleBasic, HandleResize };
  }
});
</script>

<style scoped lang="stylus"></style>
