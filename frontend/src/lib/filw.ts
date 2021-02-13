/**
 * @author 成雨
 * @date 2021/1/28
 * @Description:
 */

/**
 * 截取文件扩展名
 * @param {string} filename
 * @return {string}
 */
export function fExt(filename: string): string {
  return filename.split(".").pop() || "";
}

/**
 * 获取文件名 不带有扩展名
 * @param {string} filename
 * @returns {string}
 */
export function fName(filename: string) {
  filename = filename.replace(/\\/g, "/");
  return filename.substring(
    filename.lastIndexOf("/") + 1,
    filename.lastIndexOf(".")
  );
}

/**
 * fSize returns a pretty string from a number of bytes.
 * For example, 1024 converts to "1 MB"
 * @param {number} bytes
 * @returns {string}
 */
export function fSize(bytes: number) {
  if (bytes === 0) {
    return "0.00 B";
  }
  const e = Math.floor(Math.log(bytes) / Math.log(1024));
  return (
    (bytes / Math.pow(1024, e)).toFixed(2) + " " + " KMGTP".charAt(e) + "B"
  );
}
