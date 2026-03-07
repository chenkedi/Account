const AUTH_TOKEN_KEY = 'auth_token';
const USER_KEY = 'auth_user';

export class AuthStorage {
  async getToken(): Promise<string | null> {
    return localStorage.getItem(AUTH_TOKEN_KEY);
  }

  async setToken(token: string): Promise<void> {
    localStorage.setItem(AUTH_TOKEN_KEY, token);
  }

  async removeToken(): Promise<void> {
    localStorage.removeItem(AUTH_TOKEN_KEY);
  }

  async getUser(): Promise<{ id: string; email: string } | null> {
    const userStr = localStorage.getItem(USER_KEY);
    if (!userStr) return null;
    try {
      return JSON.parse(userStr);
    } catch {
      return null;
    }
  }

  async setUser(user: { id: string; email: string }): Promise<void> {
    localStorage.setItem(USER_KEY, JSON.stringify(user));
  }

  async removeUser(): Promise<void> {
    localStorage.removeItem(USER_KEY);
  }

  async clear(): Promise<void> {
    await this.removeToken();
    await this.removeUser();
  }
}

export const authStorage = new AuthStorage();
