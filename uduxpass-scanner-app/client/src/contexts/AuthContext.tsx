import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { scannerApi, LoginRequest } from '@/lib/api';

interface Scanner {
  id: string;
  username: string;
  name: string;
  email: string;
  role: string;
  status: string;
}

interface AuthContextType {
  scanner: Scanner | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (data: LoginRequest) => Promise<void>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [scanner, setScanner] = useState<Scanner | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    // Check for existing auth on mount
    const token = localStorage.getItem('scanner_token');
    const savedScanner = localStorage.getItem('scanner_user');
    
    if (token && savedScanner) {
      try {
        setScanner(JSON.parse(savedScanner));
      } catch (error) {
        console.error('Failed to parse saved scanner:', error);
        localStorage.removeItem('scanner_token');
        localStorage.removeItem('scanner_user');
      }
    }
    
    setIsLoading(false);
  }, []);

  const login = async (data: LoginRequest) => {
    const response = await scannerApi.login(data);
    localStorage.setItem('scanner_token', response.access_token);
    if (response.refresh_token) {
      localStorage.setItem('scanner_refresh_token', response.refresh_token);
    }
    localStorage.setItem('scanner_user', JSON.stringify(response.scanner));
    setScanner(response.scanner);
  };

  const logout = () => {
    scannerApi.logout();
    setScanner(null);
  };

  return (
    <AuthContext.Provider
      value={{
        scanner,
        isAuthenticated: !!scanner,
        isLoading,
        login,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
