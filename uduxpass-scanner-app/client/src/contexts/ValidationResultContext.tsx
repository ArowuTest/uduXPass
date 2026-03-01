import React, { createContext, useContext, useState } from 'react';

export interface ValidationResultData {
  success: boolean;
  valid: boolean;
  already_validated: boolean;
  message: string;
  serial_number?: string;
  ticket_id?: string;
  event_name?: string;
  holder_name?: string;
  ticket_tier?: string;
  validated_at?: string;
  error?: string;
}

interface ValidationResultContextType {
  result: ValidationResultData | null;
  setResult: (result: ValidationResultData | null) => void;
  clearResult: () => void;
}

const ValidationResultContext = createContext<ValidationResultContextType>({
  result: null,
  setResult: () => {},
  clearResult: () => {},
});

export function ValidationResultProvider({ children }: { children: React.ReactNode }) {
  const [result, setResult] = useState<ValidationResultData | null>(null);
  const clearResult = () => setResult(null);

  return (
    <ValidationResultContext.Provider value={{ result, setResult, clearResult }}>
      {children}
    </ValidationResultContext.Provider>
  );
}

export function useValidationResult() {
  return useContext(ValidationResultContext);
}
