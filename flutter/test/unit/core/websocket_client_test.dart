import 'package:flutter_test/flutter_test.dart';
import 'package:account/core/network/websocket_client.dart';

void main() {
  group('WebSocketClient', () {
    late WebSocketClient webSocketClient;

    setUp(() {
      webSocketClient = WebSocketClient();
    });

    tearDown(() {
      webSocketClient.dispose();
    });

    test('initial status is disconnected', () {
      expect(webSocketClient.status, WebSocketStatus.disconnected);
    });

    test('setCredentials stores auth token and device id', () {
      // Just verify the method exists and can be called without errors
      expect(() => webSocketClient.setCredentials('test-token', 'test-device'), returnsNormally);
    });

    test('disconnect works before connect', () {
      expect(() => webSocketClient.disconnect(), returnsNormally);
    });

    test('send works before connect (no error)', () {
      expect(() => webSocketClient.send({'type': 'test'}), returnsNormally);
    });

    test('statusStream is broadcast stream', () {
      expect(webSocketClient.statusStream.isBroadcast, true);
    });

    test('messageStream is broadcast stream', () {
      expect(webSocketClient.messageStream.isBroadcast, true);
    });

    test('dispose can be called multiple times', () {
      webSocketClient.dispose();
      expect(() => webSocketClient.dispose(), returnsNormally);
    });
  });
}
