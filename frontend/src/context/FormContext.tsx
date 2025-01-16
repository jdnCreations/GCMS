'use client';

import React, { createContext, useContext, useState } from 'react';

interface FormContextType {
  loginForm: boolean;
  changeFormType: () => void;
}

const FormContext = createContext<FormContextType | undefined>(undefined);

export const FormProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [loginForm, setLoginForm] = useState(true);

  const changeFormType = () => {
    setLoginForm(!loginForm);
  };

  return (
    <FormContext.Provider value={{ loginForm, changeFormType }}>
      {children}
    </FormContext.Provider>
  );
};

export const useForm = (): FormContextType => {
  const context = useContext(FormContext);
  if (!context) {
    throw new Error('useForm must be used within an FormProvider');
  }
  return context;
};
