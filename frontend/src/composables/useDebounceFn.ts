/**
 * @author 成雨
 * @date 2021/2/22
 * @Description:
 */

import { watch } from "vue";
import { useDebounce } from "./useDebounce";

/**
 * 生成带防抖功能的函数
 * @param fn
 * @param delay
 */
export function useDebounceFn<T extends (...rest: any[]) => any>(
  fn: T,
  delay = 200
) {
  const debounceValue = useDebounce(0, delay);
  let params: Parameters<T>;

  watch(
    debounceValue,
    () => {
      fn(...params);
    },
    { flush: "sync" }
  );
  return function(...rest: Parameters<T>) {
    params = rest;
    debounceValue.value++;
  };
}
