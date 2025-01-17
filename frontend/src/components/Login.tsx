import { useAuth } from '@/context/AuthContext';
import { useForm } from '@/context/FormContext';
import React, { useEffect, useState } from 'react';

const Login: React.FC = () => {
  const [errorMsg, setErrorMsg] = useState<string | null>(null);
  const { login } = useAuth();
  const [emailInput, setEmailInput] = useState('');
  const [password, setPassword] = useState('');
  const { changeFormType } = useForm();
  const { setJwt, setIsAuthenticated, setIsAdmin, setEmail, setName } =
    useAuth();
  const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

  const swapForm = () => {
    // change form to register ?
    changeFormType();
  };

  const handleLogin = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    await login({ Email: emailInput, Password: password });
  };

  return (
    <form
      className='mb-6 p-4 bg-[#B7A99A] rounded shadow text-gray-800'
      onSubmit={handleLogin}
    >
      <input
        className='mb-2 w-full p-2 border border-gray-300 focus:outline-[#a8bba1] text-[#4a4a4a] rounded'
        type='email'
        placeholder='Email'
        value={emailInput}
        onChange={(e) => {
          setEmailInput(e.target.value);
        }}
      />
      <input
        className='mb-2 w-full p-2 border border-gray-300 rounded text-[#4a4a4a] '
        type='password'
        placeholder='Password'
        value={password}
        onChange={(e) => setPassword(e.target.value)}
      />
      <button
        type='submit'
        className='bg-[#A8BBA1] text-white w-full p-2 rounded hover:bg-[#C8D9C3]'
      >
        Login
      </button>
      <p className='text-[#4a4a4a] pt-2'>
        Not a member yet?{' '}
        <button className='underline' type='button' onClick={swapForm}>
          Register now
        </button>
      </p>
      {errorMsg && <p className='text-red-500 font-bold'>{errorMsg}</p>}
    </form>
  );
};

export default Login;
