import { onMounted, ref, Ref } from "vue";

// 处理文件拖拽到浏览器的相关操作
export default function useDrag<T extends HTMLElement | undefined>(
  elRef: Ref<T>
) {
  const files = ref<FileList>();

  onMounted(() => {
    if (elRef.value) {
      elRef.value.addEventListener(
        "dragenter",
        e => {
          e.stopPropagation();
          e.preventDefault();
        },
        false
      );
      elRef.value.addEventListener(
        "dragover",
        e => {
          e.stopPropagation();
          e.preventDefault();
        },
        false
      );
      elRef.value.addEventListener(
        "dragend",
        e => {
          e.stopPropagation();
          e.preventDefault();
        },
        false
      );
      elRef.value.addEventListener(
        "dragleave",
        e => {
          e.stopPropagation();
          e.preventDefault();
        },
        false
      );
      elRef.value.addEventListener(
        "drop",
        e => {
          e.stopPropagation();
          e.preventDefault();
          if (e.dataTransfer) files.value = e.dataTransfer.files;
        },
        false
      );
    }
  });

  return {
    files
  };
}
