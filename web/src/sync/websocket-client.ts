import ReconnectingWebSocket from 'reconnecting-websocket';
import { Subject, BehaviorSubject } from 'rxjs';
import { API_CONSTANTS } from '../core/constants/api.constants';

export type WebSocketStatus = 'disconnected' | 'connecting' | 'connected' | 'error';

export class WebSocketClient {
  private ws: ReconnectingWebSocket | null = null;
  private statusSubject = new BehaviorSubject<WebSocketStatus>('disconnected');
  private messageSubject = new Subject<unknown>();

  public status$ = this.statusSubject.asObservable();
  public message$ = this.messageSubject.asObservable();

  constructor(
    private authToken: string,
    private deviceId: string
  ) {}

  connect(): void {
    if (this.ws) return;

    this.updateStatus('connecting');

    const wsUrl = new URL(
      `${API_CONSTANTS.baseUrl.replace('http', 'ws')}${API_CONSTANTS.wsSync}`
    );
    wsUrl.searchParams.set('token', this.authToken);
    wsUrl.searchParams.set('device_id', this.deviceId);

    this.ws = new ReconnectingWebSocket(wsUrl.toString());

    this.ws.onopen = () => {
      this.updateStatus('connected');
    };

    this.ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        this.messageSubject.next(data);
      } catch (e) {
        console.error('WebSocket message parse error:', e);
      }
    };

    this.ws.onerror = () => {
      this.updateStatus('error');
    };

    this.ws.onclose = () => {
      this.updateStatus('disconnected');
    };
  }

  disconnect(): void {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }

  send(message: unknown): void {
    if (this.ws && this.statusSubject.value === 'connected') {
      this.ws.send(JSON.stringify(message));
    }
  }

  private updateStatus(status: WebSocketStatus): void {
    this.statusSubject.next(status);
  }

  dispose(): void {
    this.disconnect();
    this.statusSubject.complete();
    this.messageSubject.complete();
  }
}
