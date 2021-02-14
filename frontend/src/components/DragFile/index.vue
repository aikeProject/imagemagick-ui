<template>
  <div
    class="flex flex-1 m-2 box-border border-2 rounded-md border-red-500 border-dashed"
  >
    <label class="flex-1" :ref="bindDragRef">
      <input hidden type="file" />
    </label>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, watch } from "vue";
import useDrag from "composables/useDrag";

export default defineComponent({
  name: "Home",
  components: {},
  setup() {
    const dragRef = ref<HTMLLabelElement>();
    const { files } = useDrag<HTMLLabelElement | undefined>(dragRef);
    watch(files, () => {
      console.log(files.value);
    });
    const bindDragRef = (el: HTMLLabelElement) => {
      dragRef.value = el;
    };

    return { bindDragRef };
  }
});
</script>

<style scoped lang="stylus"></style>
