import 'dart:async';
import 'dart:convert';

import 'package:web_socket_channel/web_socket_channel.dart';

import '../constants/api_constants.dart';

enum WebSocketStatus { disconnected, connecting, connected, error }

class WebSocketClient {
  WebSocketChannel? _channel;
  final _statusController = StreamController<WebSocketStatus>.broadcast();
  final _messageController = StreamController<Map<String, dynamic>>.broadcast();

  WebSocketStatus _status = WebSocketStatus.disconnected;
  String? _authToken;
  String? _deviceId;
  Timer? _reconnectTimer;

  Stream<WebSocketStatus> get statusStream => _statusController.stream;
  Stream<Map<String, dynamic>> get messageStream => _messageController.stream;
  WebSocketStatus get status => _status;

  void setCredentials(String authToken, String deviceId) {
    _authToken = authToken;
    _deviceId = deviceId;
  }

  void connect() {
    if (_authToken == null || _deviceId == null) {
      return;
    }

    _updateStatus(WebSocketStatus.connecting);

    try {
      final wsUrl = Uri.parse(
        '${ApiConstants.baseUrl.replaceFirst('http', 'ws')}${ApiConstants.wsSync}',
      ).replace(
        queryParameters: {
          'token': _authToken,
          'device_id': _deviceId,
        },
      );

      _channel = WebSocketChannel.connect(wsUrl);
      _updateStatus(WebSocketStatus.connected);

      _channel!.stream.listen(
        (message) {
          try {
            final data = jsonDecode(message as String) as Map<String, dynamic>;
            _messageController.add(data);
          } catch (e) {
            // Ignore invalid messages
          }
        },
        onError: (_) {
          _updateStatus(WebSocketStatus.error);
          _scheduleReconnect();
        },
        onDone: () {
          _updateStatus(WebSocketStatus.disconnected);
          _scheduleReconnect();
        },
      );
    } catch (e) {
      _updateStatus(WebSocketStatus.error);
      _scheduleReconnect();
    }
  }

  void disconnect() {
    _reconnectTimer?.cancel();
    _reconnectTimer = null;
    _channel?.sink.close();
    _channel = null;
    _updateStatus(WebSocketStatus.disconnected);
  }

  void send(Map<String, dynamic> message) {
    if (_status == WebSocketStatus.connected && _channel != null) {
      _channel!.sink.add(jsonEncode(message));
    }
  }

  void _updateStatus(WebSocketStatus status) {
    _status = status;
    _statusController.add(status);
  }

  void _scheduleReconnect() {
    if (_reconnectTimer != null) return;

    _reconnectTimer = Timer(const Duration(seconds: 5), () {
      _reconnectTimer = null;
      if (_status != WebSocketStatus.connected) {
        connect();
      }
    });
  }

  void dispose() {
    disconnect();
    _statusController.close();
    _messageController.close();
  }
}
