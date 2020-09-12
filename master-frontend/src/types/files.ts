export interface Node {
  path: string,
  info: FileInfo,
  children: Node[],
  parent_name: string,
}

export interface FileInfo {
  name: string,
  size: number,
  is_dir: boolean,
  content: string
}
