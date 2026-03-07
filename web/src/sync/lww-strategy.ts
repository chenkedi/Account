import type { LwwEntity } from '../data/models/base.model';

export class LwwStrategy {
  /**
   * 比较两个实体，返回是否应该保留本地实体
   * 如果本地更新时间更新，返回 true；否则返回 false
   */
  static shouldKeepLocal<T extends LwwEntity>(local: T, remote: T): boolean {
    return new Date(local.last_modified_at) > new Date(remote.last_modified_at);
  }

  /**
   * 合并两个实体列表，使用 LWW 策略解决冲突
   */
  static mergeLists<T extends LwwEntity>(local: T[], remote: T[]): T[] {
    const merged = new Map<string, T>();

    // 先添加所有本地实体
    for (const entity of local) {
      merged.set(entity.id, entity);
    }

    // 合并远程实体
    for (const remoteEntity of remote) {
      const existing = merged.get(remoteEntity.id);
      if (!existing) {
        merged.set(remoteEntity.id, remoteEntity);
      } else {
        if (!this.shouldKeepLocal(existing, remoteEntity)) {
          merged.set(remoteEntity.id, remoteEntity);
        }
      }
    }

    return Array.from(merged.values());
  }

  /**
   * 解析单个实体的冲突
   */
  static resolve<T extends LwwEntity>(local: T | undefined, remote: T): T {
    if (!local) return remote;
    return this.shouldKeepLocal(local, remote) ? local : remote;
  }
}
