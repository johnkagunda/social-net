'use client';

import { createContext, useContext, useState, useEffect } from 'react';
import { getMe, login as loginApi, logout as logoutApi, getSession } from '@/lib/auth';

const AuthContext = createContext();

export { AuthContext };

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [sessionToken, setSessionToken] = useState(null);

  useEffect(() => {
    checkAuth();
  }, []);

  const checkAuth = async () => {
    try {
      const userData = await getMe();
      setUser(userData);
      if (userData) {
        // Get session token for WebSocket - non-fatal
        try {
          const sessionData = await getSession();
          if (sessionData && sessionData.session_id) {
            setSessionToken(sessionData.session_id);
          }
        } catch (e) {}
      }
    } catch (error) {
      setUser(null);
      setSessionToken(null);
    } finally {
      setLoading(false);
    }
  };

  const login = async (email, password) => {
    const userData = await loginApi(email, password);
    setUser(userData);
    // Get session token for WebSocket after login - non-fatal if it fails
    try {
      const sessionData = await getSession();
      if (sessionData && sessionData.session_id) {
        setSessionToken(sessionData.session_id);
      }
    } catch (e) {
      // WebSocket session is optional - login still succeeds
    }
    return userData;
  };

  const logout = async () => {
    await logoutApi();
    setUser(null);
    setSessionToken(null);
  };

  return (
    <AuthContext.Provider value={{ user, loading, login, logout, checkAuth, sessionToken }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
