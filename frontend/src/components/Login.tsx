import { useAuth } from '@/context/AuthContext';
import React, { useState } from 'react';

const Login: React.FC<LoginProps> = () => {
  const [errorMsg, setErrorMsg] = useState<string | null>(null);
  const { isAuthenticated, login, logout, email, jwt } = useAuth();
  const [emailInput, setEmailInput] = useState('');
  const [password, setPassword] = useState('');

  const handleLogin = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    await login({ Email: emailInput, Password: password });
  };

  return (
    <form
      className='mb-6 p-4 bg-blue-400 rounded shadow flex flex-col w-full'
      onSubmit={handleLogin}
    >
      <input
        className='mb-2 w-full p-2 border border-gray-300 rounded'
        type='email'
        placeholder='Email'
        value={emailInput}
        onChange={(e) => {
          setEmailInput(e.target.value);
        }}
      />
      <input
        className='mb-2 w-full p-2 border border-gray-300 rounded'
        type='password'
        placeholder='Password'
        value={password}
        onChange={(e) => setPassword(e.target.value)}
      />
      <button
        type='submit'
        className='bg-cyan-400 w-full p-2 rounded hover:bg-cyan-500'
      >
        Login
      </button>
      {errorMsg && <p className='text-red-500 font-bold'>{errorMsg}</p>}
    </form>
  );
};

export default Login;
