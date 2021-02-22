/**
 * @author 成雨
 * @date 2021/2/22
 * @Description:
 */

import { customRef, getCurrentInstance, onUnmounted } from "vue";

/**
 * 带防抖功能的状态
 * @param value
 * @param delay
 */
export function useDebounce<T>(value: T, delay = 200) {
  let timer: any;
  const clear = () => {
    if (timer) {
      clearTimeout(timer);
    }
  };
  if (getCurrentInstance()) {
    onUnmounted(() => {
      clear();
    });
  }
  return customRef((tracker, trigger) => ({
    get() {
      tracker();
      return value;
    },
    set(val: T) {
      clear();
      timer = setTimeout(() => {
        value = val;
        timer = null;
        trigger();
      }, delay);
    }
  }));
}
