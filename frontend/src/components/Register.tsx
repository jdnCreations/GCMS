import { useAuth } from '@/context/AuthContext';
import { useForm } from '@/context/FormContext';
import { useState } from 'react';

const Register: React.FC = () => {
  const { changeFormType } = useForm();
  const { error, register, setError } = useAuth();

  const swapForm = () => {
    changeFormType();
    setError(null);
  };

  const handleRegister = async (e: React.FormEvent<HTMLFormElement>) => {
    await register(newUser, e);
  };

  const [newUser, setNewUser] = useState({
    FirstName: '',
    LastName: '',
    Email: '',
    Password: '',
  });

  /* Create user */
  return (
    <form
      onSubmit={handleRegister}
      className='mb-6 p-4 bg-[#B7A99A] rounded shadow text-gray-800'
    >
      <input
        type='text'
        placeholder='First Name'
        value={newUser.FirstName}
        onChange={(e) => setNewUser({ ...newUser, FirstName: e.target.value })}
        className='mb-2 w-full p-2 border text-[#4a4a4a] border-gray-300 rounded focus:outline-nook-olive'
      />
      <input
        type='text'
        placeholder='Last Name'
        value={newUser.LastName}
        onChange={(e) => setNewUser({ ...newUser, LastName: e.target.value })}
        className='mb-2 w-full p-2 border text-[#4a4a4a] border-gray-300 rounded focus:outline-nook-olive'
      />
      <input
        type='text'
        placeholder='Email'
        value={newUser.Email}
        onChange={(e) => setNewUser({ ...newUser, Email: e.target.value })}
        className='mb-2 w-full p-2 border text-[#4a4a4a] border-gray-300 rounded focus:outline-nook-olive'
      />
      <input
        type='password'
        placeholder='Password'
        value={newUser.Password}
        onChange={(e) => setNewUser({ ...newUser, Password: e.target.value })}
        className='mb-2 w-full p-2 border text-[#4a4a4a] border-gray-300 rounded focus:outline-nook-olive'
      />
      <button
        className='bg-nook-olive text-white w-full p-2 rounded hover:bg-nook-light-olive'
        type='submit'
      >
        Register
      </button>
      <p className='text-nook-charcoal pt-2'>
        Already a member?{' '}
        <button className='underline' type='button' onClick={swapForm}>
          Login
        </button>
        {error && <p className='text-red-500 font-bold'>{error}</p>}
      </p>
    </form>
  );
};

export default Register;
