/// Last-Write-Wins conflict resolution strategy
/// Compares entities based on their last_modified_at timestamp
class LwwStrategy {
  /// Determine which version of an entity should win
  /// Returns true if local should be kept, false if remote should be used
  static bool shouldKeepLocal<T extends LwwEntity>(T local, T remote) {
    return local.lastModifiedAt.isAfter(remote.lastModifiedAt);
  }

  /// Merge two lists of entities using LWW strategy
  static List<T> mergeLists<T extends LwwEntity>(List<T> local, List<T> remote) {
    final Map<String, T> merged = {};

    // Add all local entities
    for (final entity in local) {
      merged[entity.id] = entity;
    }

    // Merge with remote entities
    for (final remoteEntity in remote) {
      final existing = merged[remoteEntity.id];
      if (existing == null) {
        merged[remoteEntity.id] = remoteEntity;
      } else {
        if (!shouldKeepLocal(existing, remoteEntity)) {
          merged[remoteEntity.id] = remoteEntity;
        }
      }
    }

    return merged.values.toList();
  }
}

abstract class LwwEntity {
  String get id;
  DateTime get lastModifiedAt;
  bool get isDeleted;
  int get version;
}
