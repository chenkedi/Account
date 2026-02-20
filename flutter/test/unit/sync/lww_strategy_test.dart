import 'package:flutter_test/flutter_test.dart';
import 'package:account/sync/lww_strategy.dart';
import '../../helpers/test_data.dart';

void main() {
  group('LwwStrategy', () {
    group('shouldKeepLocal', () {
      test('returns true when local entity is newer', () {
        final now = DateTime.now().toUtc();
        final local = TestData.createTransaction(
          lastModifiedAt: now,
        );
        final remote = TestData.createTransaction(
          id: local.id,
          lastModifiedAt: now.subtract(const Duration(hours: 1)),
        );

        expect(LwwStrategy.shouldKeepLocal(local, remote), true);
      });

      test('returns false when remote entity is newer', () {
        final now = DateTime.now().toUtc();
        final local = TestData.createTransaction(
          lastModifiedAt: now.subtract(const Duration(hours: 1)),
        );
        final remote = TestData.createTransaction(
          id: local.id,
          lastModifiedAt: now,
        );

        expect(LwwStrategy.shouldKeepLocal(local, remote), false);
      });

      test('returns false when timestamps are equal', () {
        final now = DateTime.now().toUtc();
        final local = TestData.createTransaction(
          lastModifiedAt: now,
        );
        final remote = TestData.createTransaction(
          id: local.id,
          lastModifiedAt: now,
        );

        expect(LwwStrategy.shouldKeepLocal(local, remote), false);
      });
    });

    group('mergeLists', () {
      test('combines entities with different IDs', () {
        final local1 = TestData.createTransaction(id: 'local-1');
        final local2 = TestData.createTransaction(id: 'local-2');
        final remote1 = TestData.createTransaction(id: 'remote-1');
        final remote2 = TestData.createTransaction(id: 'remote-2');

        final merged = LwwStrategy.mergeLists(
          [local1, local2],
          [remote1, remote2],
        );

        expect(merged.length, 4);
        expect(merged.any((t) => t.id == 'local-1'), true);
        expect(merged.any((t) => t.id == 'local-2'), true);
        expect(merged.any((t) => t.id == 'remote-1'), true);
        expect(merged.any((t) => t.id == 'remote-2'), true);
      });

      test('uses local entity when it is newer', () {
        final now = DateTime.now().toUtc();
        final local = TestData.createTransaction(
          id: 'conflict-1',
          note: 'Local version',
          lastModifiedAt: now,
        );
        final remote = TestData.createTransaction(
          id: 'conflict-1',
          note: 'Remote version',
          lastModifiedAt: now.subtract(const Duration(hours: 1)),
        );

        final merged = LwwStrategy.mergeLists([local], [remote]);

        expect(merged.length, 1);
        expect(merged.first.note, 'Local version');
      });

      test('uses remote entity when it is newer', () {
        final now = DateTime.now().toUtc();
        final local = TestData.createTransaction(
          id: 'conflict-1',
          note: 'Local version',
          lastModifiedAt: now.subtract(const Duration(hours: 1)),
        );
        final remote = TestData.createTransaction(
          id: 'conflict-1',
          note: 'Remote version',
          lastModifiedAt: now,
        );

        final merged = LwwStrategy.mergeLists([local], [remote]);

        expect(merged.length, 1);
        expect(merged.first.note, 'Remote version');
      });

      test('uses remote entity when timestamps are equal', () {
        final now = DateTime.now().toUtc();
        final local = TestData.createTransaction(
          id: 'conflict-1',
          note: 'Local version',
          lastModifiedAt: now,
        );
        final remote = TestData.createTransaction(
          id: 'conflict-1',
          note: 'Remote version',
          lastModifiedAt: now,
        );

        final merged = LwwStrategy.mergeLists([local], [remote]);

        expect(merged.length, 1);
        expect(merged.first.note, 'Remote version');
      });
    });
  });
}
