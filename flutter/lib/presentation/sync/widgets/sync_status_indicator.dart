import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import '../bloc/sync_bloc.dart';
import '../../../sync/sync_manager.dart' as sync_manager;
import '../../../core/utils/date_utils.dart' as date_utils;

class SyncStatusIndicator extends StatelessWidget {
  const SyncStatusIndicator({super.key});

  @override
  Widget build(BuildContext context) {
    return BlocBuilder<SyncBloc, SyncState>(
      builder: (context, state) {
        return IconButton(
          icon: _buildIcon(context, state),
          onPressed: state.isSyncing
              ? null
              : () {
                  context.read<SyncBloc>().add(const SyncRequested());
                },
          tooltip: _getTooltip(state),
        );
      },
    );
  }

  Widget _buildIcon(BuildContext context, SyncState state) {
    switch (state.status) {
      case SyncStatus.syncing:
        return SizedBox(
          width: 24,
          height: 24,
          child: CircularProgressIndicator(
            strokeWidth: 2,
            valueColor: AlwaysStoppedAnimation<Color>(
              Theme.of(context).colorScheme.onSurface,
            ),
          ),
        );
      case SyncStatus.success:
        return Icon(
          Icons.cloud_done,
          color: Theme.of(context).colorScheme.primary,
        );
      case SyncStatus.error:
        return Icon(
          Icons.cloud_off,
          color: Theme.of(context).colorScheme.error,
        );
      case SyncStatus.idle:
      default:
        return Icon(
          Icons.cloud_sync,
          color: Theme.of(context).colorScheme.onSurface,
        );
    }
  }

  String _getTooltip(SyncState state) {
    switch (state.status) {
      case SyncStatus.syncing:
        return state.message ?? 'Syncing...';
      case SyncStatus.success:
        if (state.lastSyncAt != null) {
          return 'Last synced: ${date_utils.DateUtils.formatDateTime(state.lastSyncAt!)}';
        }
        return 'Sync completed';
      case SyncStatus.error:
        return state.errorMessage ?? 'Sync failed';
      case SyncStatus.idle:
      default:
        if (state.lastSyncAt != null) {
          return 'Last synced: ${date_utils.DateUtils.formatDateTime(state.lastSyncAt!)}';
        }
        return 'Tap to sync';
    }
  }
}
