import { useAuth } from '@/context/AuthContext';
import { useForm } from '@/context/FormContext';
import React, { useState } from 'react';

const Login: React.FC = () => {
  const { login, error, setError } = useAuth();
  const [emailInput, setEmailInput] = useState('');
  const [password, setPassword] = useState('');
  const { changeFormType } = useForm();

  const swapForm = () => {
    changeFormType();
    setError(null);
  };

  const handleLogin = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    await login({ Email: emailInput, Password: password }, e);
  };

  return (
    <form
      className='mb-6 p-4 bg-[#B7A99A] rounded shadow text-gray-800'
      onSubmit={handleLogin}
    >
      <input
        className='mb-2 w-full p-2 border border-gray-300 focus:outline-nook-olive text-[#4a4a4a] rounded'
        type='email'
        placeholder='Email'
        value={emailInput}
        onChange={(e) => {
          setEmailInput(e.target.value);
        }}
      />
      <input
        className='mb-2 w-full p-2 border border-gray-300 rounded text-[#4a4a4a] focus:outline-nook-olive '
        type='password'
        placeholder='Password'
        value={password}
        onChange={(e) => setPassword(e.target.value)}
      />
      <button
        type='submit'
        className='bg-nook-olive text-white w-full p-2 rounded hover:bg-nook-light-olive'
      >
        Log in
      </button>
      <p className='text-[#4a4a4a] pt-2'>
        Not a member yet?{' '}
        <button className='underline' type='button' onClick={swapForm}>
          Register now
        </button>
      </p>
      {error && <p className='text-red-500 font-bold'>{error}</p>}
    </form>
  );
};

export default Login;
