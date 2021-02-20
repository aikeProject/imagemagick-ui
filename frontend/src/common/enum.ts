/**
 * @author 成雨
 * @date 2021/2/19
 * @Description:
 */

export type FileStatusValue = 0 | 1 | 2 | 3 | 4;

// 文件状态
export const enum FileStatus {
  NotStarted,
  Start,
  SendSuccess,
  Running,
  Done
}
