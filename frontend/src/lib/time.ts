/**
 * @author 成雨
 * @date 2021/1/31
 * @Description:
 */

/**
 * 毫秒时间戳格式化
 * @param {number} ms
 * @returns {(string)[]}
 */
export function prettyTime(ms: number) {
  const seconds = (ms / 1000).toFixed(1);
  const minutes = (ms / (1000 * 60)).toFixed(1);
  const hours = (ms / (1000 * 60 * 60)).toFixed(1);
  const days = (ms / (1000 * 60 * 60 * 24)).toFixed(1);

  if (ms < 1000) {
    return [ms, "毫秒"];
  } else if (Number(seconds) < 60) {
    return [seconds, "秒"];
  } else if (Number(minutes) < 60) {
    return [minutes, "分"];
  } else if (Number(hours) < 24) {
    return [hours, "小时"];
  } else {
    return [days, "天"];
  }
}
