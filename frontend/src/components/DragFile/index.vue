<template>
  <label
    :ref="bindDragRef"
    class="m-2 flex flex-1 justify-center items-center px-6 pt-5 pb-6 border-2 border-gray-300 border-dashed rounded-md hover:border-blue-500"
  >
    <div class="space-y-1 text-center">
      <svg
        class="mx-auto h-12 w-12 text-gray-400"
        stroke="currentColor"
        fill="none"
        viewBox="0 0 48 48"
        aria-hidden="true"
      >
        <path
          d="M28 8H12a4 4 0 00-4 4v20m32-12v8m0 0v8a4 4 0 01-4 4H12a4 4 0 01-4-4v-4m32-4l-3.172-3.172a4 4 0 00-5.656 0L28 28M8 32l9.172-9.172a4 4 0 015.656 0L28 28m0 0l4 4m4-24h8m-4-4v8m-12 4h.02"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        />
      </svg>
      <div class="flex text-sm text-gray-600">
        <div
          class="relative cursor-pointer bg-white rounded-md font-medium text-indigo-600 hover:text-indigo-500"
        >
          <span>file</span>
          <input
            hidden
            id="file-upload"
            name="file-upload"
            type="file"
            class="sr-only"
          />
        </div>
        <p class="pl-1">or drag and drop</p>
      </div>
      <p class="text-xs text-gray-500">
        PNG, JPG, GIF
      </p>
    </div>
  </label>
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
