"use client";

import { useEffect, useState } from "react";
import { getAuthState, subscribeAuth } from "./authStore";

export function useAuth() {
  const [auth, setAuth] = useState(getAuthState());

  useEffect(() => {
    return subscribeAuth(setAuth);
  }, []);

  return auth; // { isAuthenticated }
}
