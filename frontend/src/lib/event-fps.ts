/**
 * @author 成雨
 * @date 2021/2/20
 * @Description:
 */

import mitt from "mitt";
export const EventFps = mitt();
// 计算性能指标
(() => {
  const createConsole = (desc: string, val: number) =>
    console.log(
      `%c${desc}`,
      "color:#fff;background:red;padding:2px 6px;border-radius:3px;",
      val
    );

  window.addEventListener("load", () => {
    const timing = performance.timing;
    createConsole(
      "DNS 解析耗时",
      timing.domainLookupEnd - timing.domainLookupStart
    );
    createConsole("TCP连接耗时", timing.connectEnd - timing.connectStart);
    createConsole("网络请求耗时", timing.responseStart - timing.requestStart);
    createConsole("数据传输耗时", timing.responseEnd - timing.requestStart);
    createConsole(
      "页面首次渲染时间",
      timing.responseEnd - timing.navigationStart
    );
    createConsole(
      "首次可交互时间",
      timing.domInteractive - timing.navigationStart
    );
    createConsole("DOM解析耗时", timing.domInteractive - timing.responseEnd);
    createConsole("DOM构建耗时", timing.domComplete - timing.domInteractive);
    createConsole(
      "HTML 加载完成,DOM Ready",
      timing.domContentLoadedEventEnd - timing.navigationStart
    );
    createConsole(
      "页面完全加载耗时",
      timing.loadEventStart - timing.navigationStart
    );
  });
})();

let start = 0;
(function loop(time?: number) {
  time = time || 0;
  // 单位 帧/s
  const duration = (time - start) / 1000;
  start = time;
  EventFps.emit("update", duration);
  requestAnimationFrame(loop);
})();
