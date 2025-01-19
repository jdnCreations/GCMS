import axios from 'axios';
import { useState } from 'react';
import MessageDisplay from './MessageDisplay';
import { useAuth } from '@/context/AuthContext';

interface UpdatedUser {
  FirstName?: string;
  LastName?: string;
  Email?: string;
}

const UpdateUser: React.FC = () => {
  const { jwt, userId } = useAuth();
  const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
  const [updateMsg, setUpdateMsg] = useState('');
  const [errMsg, setErrMsg] = useState('');
  const [updatedUserInfo, setUpdatedUserInfo] = useState<UpdatedUser>({
    FirstName: '',
    LastName: '',
    Email: '',
  });
  const updateUser = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    try {
      await axios.put(`${apiUrl}/api/users/${userId}`, updatedUserInfo, {
        headers: {
          Authorization: `Bearer ${jwt}`,
        },
      });
      setUpdatedUserInfo({ FirstName: '', LastName: '', Email: '' });
      // tell user their info was updated successfully
      setUpdateMsg('updated successfully');
    } catch (error) {
      if (axios.isAxiosError(error) && error.response) {
        setErrMsg(error?.response?.data.error);
      }
      console.log(error);
    }
  };
  return (
    <div>
      <form
        onSubmit={updateUser}
        className='mb-6 p-4 bg-[#B7A99A] rounded shadow text-gray-800'
      >
        <h1 className='font-bold text-2xl'>Update User Details</h1>
        <MessageDisplay errorMsg={errMsg} updateMsg={updateMsg} />
        <input
          type='text'
          placeholder='First Name'
          value={updatedUserInfo.FirstName}
          onChange={(e) =>
            setUpdatedUserInfo({
              ...updatedUserInfo,
              FirstName: e.target.value,
            })
          }
          className='mb-2 w-full p-2 border text-[#4a4a4a] border-gray-300 rounded focus:outline-nook-olive'
        />
        <input
          type='text'
          placeholder='Last Name'
          value={updatedUserInfo.LastName}
          onChange={(e) =>
            setUpdatedUserInfo({
              ...updatedUserInfo,
              LastName: e.target.value,
            })
          }
          className='mb-2 w-full p-2 border text-[#4a4a4a] border-gray-300 rounded focus:outline-nook-olive'
        />
        <input
          type='text'
          placeholder='Email'
          value={updatedUserInfo.Email}
          onChange={(e) =>
            setUpdatedUserInfo({ ...updatedUserInfo, Email: e.target.value })
          }
          className='mb-2 w-full p-2 border text-[#4a4a4a] border-gray-300 rounded focus:outline-nook-olive '
        />
        <button
          className='bg-nook-olive text-white w-full p-2 rounded hover:bg-nook-light-olive'
          type='submit'
        >
          Update Details
        </button>
      </form>
    </div>
  );
};

export default UpdateUser;
