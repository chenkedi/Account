type EventCallback = (...args: unknown[]) => void;

interface EventEmitter {
  on(event: string, callback: EventCallback): void;
  off(event: string, callback: EventCallback): void;
  emit(event: string, ...args: unknown[]): void;
}

class AuthEventEmitter implements EventEmitter {
  private events: Map<string, EventCallback[]> = new Map();

  on(event: string, callback: EventCallback): void {
    if (!this.events.has(event)) {
      this.events.set(event, []);
    }
    this.events.get(event)!.push(callback);
  }

  off(event: string, callback: EventCallback): void {
    const callbacks = this.events.get(event);
    if (callbacks) {
      this.events.set(
        event,
        callbacks.filter((cb) => cb !== callback)
      );
    }
  }

  emit(event: string, ...args: unknown[]): void {
    const callbacks = this.events.get(event);
    if (callbacks) {
      callbacks.forEach((callback) => {
        try {
          callback(...args);
        } catch (error) {
          console.error('[AuthEventEmitter] Error in callback:', error);
        }
      });
    }
  }
}

// 导出单例实例
export const authEvents = new AuthEventEmitter();

// 导出事件类型常量
export const AUTH_EVENTS = {
  SESSION_EXPIRED: 'auth:session-expired',
} as const;
