type AuthState = {
  isAuthenticated: boolean;
};

let state: AuthState = { isAuthenticated: false };
const listeners = new Set<(s: AuthState) => void>();

export function getAuthState() {
  return state;
}

export function setAuthenticated(isAuthenticated: boolean) {
  state = { ...state, isAuthenticated };
  listeners.forEach((l) => l(state));
}

export function subscribeAuth(listener: (s: AuthState) => void) {
  listeners.add(listener);
  return () => {
    listeners.delete(listener); // 👈 ne returnuj boolean
  };
}
