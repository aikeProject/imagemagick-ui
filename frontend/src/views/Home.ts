/**
 * @author 成雨
 * @date 2021/2/19
 * @Description:
 */

import { FileStatusValue } from "common/enum";

export interface FileData {
  name: string;
  size: number;
  src: string;
  status: FileStatusValue;
}
