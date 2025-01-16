import { useForm } from '@/context/FormContext';
import axios from 'axios';
import { useState } from 'react';

const Register: React.FC = () => {
  const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
  const { changeFormType } = useForm();

  const swapForm = () => changeFormType();

  const [newUser, setNewUser] = useState({
    FirstName: '',
    LastName: '',
    Email: '',
    Password: '',
  });

  const createUser = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    try {
      await axios.post(`${apiUrl}/api/users`, newUser);
      setNewUser({ FirstName: '', LastName: '', Email: '', Password: '' });
    } catch (error) {
      if (axios.isAxiosError(error) && error.response) {
        console.error('API Error:', error.response.data.error);
      } else {
        console.error('Error creating user:', error);
      }
    }
  };
  /* Create user */
  return (
    <form
      onSubmit={createUser}
      className='mb-6 p-4 bg-[#B7A99A] rounded shadow text-gray-800'
    >
      <input
        type='text'
        placeholder='First Name'
        value={newUser.FirstName}
        onChange={(e) => setNewUser({ ...newUser, FirstName: e.target.value })}
        className='mb-2 w-full p-2 border text-[#4a4a4a] border-gray-300 rounded'
      />
      <input
        type='text'
        placeholder='Last Name'
        value={newUser.LastName}
        onChange={(e) => setNewUser({ ...newUser, LastName: e.target.value })}
        className='mb-2 w-full p-2 border text-[#4a4a4a] border-gray-300 rounded'
      />
      <input
        type='text'
        placeholder='Email'
        value={newUser.Email}
        onChange={(e) => setNewUser({ ...newUser, Email: e.target.value })}
        className='mb-2 w-full p-2 border text-[#4a4a4a] border-gray-300 rounded'
      />
      <input
        type='password'
        placeholder='Password'
        value={newUser.Password}
        onChange={(e) => setNewUser({ ...newUser, Password: e.target.value })}
        className='mb-2 w-full p-2 border text-[#4a4a4a] border-gray-300 rounded'
      />
      <button
        className='bg-[#A8BBA1] text-white w-full p-2 rounded hover:bg-[#E4B7B2]'
        type='submit'
      >
        Register
      </button>
      <p className='text-[#4a4a4a] pt-2'>
        Already a member?{' '}
        <button className='underline' type='button' onClick={swapForm}>
          Login
        </button>
      </p>
    </form>
  );
};

export default Register;
