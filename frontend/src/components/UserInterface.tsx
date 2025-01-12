'use client';

import React, { useEffect, useState } from 'react';
import CardComponent from './CardComponent';
import axios from 'axios';

interface User {
  ID: string;
  FirstName: string;
  LastName: string;
  Email: string;
}

interface UserInterfaceProps {
  backendName: string;
}

const UserInterface: React.FC<UserInterfaceProps> = ({ backendName }) => {
  const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
  const [users, setUsers] = useState<User[]>([]);
  const [newUser, setNewUser] = useState({
    FirstName: '',
    LastName: '',
    Email: '',
  });
  const [updatedUser, setUpdateUser] = useState({
    ID: '',
    FirstName: '',
    LastName: '',
    Email: '',
  });

  // fetch all users
  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.get(`${apiUrl}/api/users`);
        if (response.data && response.data.length > 0) {
          setUsers(response?.data.reverse());
        }
      } catch (error) {
        console.error('Error fetching data:', error);
      }
    };

    fetchData();
  }, [backendName, apiUrl]);

  const createUser = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    try {
      const response = await axios.post(`${apiUrl}/api/users`, newUser);
      setUsers([response.data, ...users]);
      setNewUser({ FirstName: '', LastName: '', Email: '' });
    } catch (error) {
      if (axios.isAxiosError(error) && error.response) {
        console.error('API Error:', error.response.data.error);
      } else {
        console.error('Error creating user:', error);
      }
    }
  };

  const updateUser = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    try {
      const response = await axios.put(
        `${apiUrl}/api/users/${updatedUser.ID}`,
        updatedUser
      );
      const updatedUserData = response.data;
      setUsers(
        users.map((user) =>
          user.ID === updatedUserData.ID ? updatedUserData : user
        )
      );
      setUpdateUser({ ID: '', FirstName: '', LastName: '', Email: '' });
    } catch (error) {
      if (axios.isAxiosError(error) && error.response) {
        console.error('API Error:', error.response.data.error);
      } else {
        console.error('Error updating user:', error);
      }
    }
  };

  const deleteUser = async (ID: string) => {
    try {
      await axios.delete(`${apiUrl}/api/users/${ID}`);
      setUsers(users.filter((user) => user.ID !== ID));
    } catch (error) {
      console.error('Error deleting user:', error);
    }
  };

  return (
    <div className='bg-blue-100 justify-center items-center flex flex-col rounded-lg py-2 px-8'>
      {/* Create user */}
      <div>
        <form
          onSubmit={createUser}
          className='mb-6 p-4 bg-blue-400 rounded shadow'
        >
          <input
            type='text'
            placeholder='First Name'
            value={newUser.FirstName}
            onChange={(e) =>
              setNewUser({ ...newUser, FirstName: e.target.value })
            }
            className='mb-2 w-full p-2 border border-gray-300 rounded'
          />
          <input
            type='text'
            placeholder='Last Name'
            value={newUser.LastName}
            onChange={(e) =>
              setNewUser({ ...newUser, LastName: e.target.value })
            }
            className='mb-2 w-full p-2 border border-gray-300 rounded'
          />
          <input
            type='text'
            placeholder='Email'
            value={newUser.Email}
            onChange={(e) => setNewUser({ ...newUser, Email: e.target.value })}
            className='mb-2 w-full p-2 border border-gray-300 rounded'
          />
          <button
            className='bg-cyan-400 w-full p-2 rounded hover:bg-cyan-500'
            type='submit'
          >
            Create User
          </button>
        </form>
      </div>

      {/* Update user */}
      <div>
        <form
          onSubmit={updateUser}
          className='mb-6 p-4 bg-blue-100 rounded shadow'
        >
          <input
            type='text'
            placeholder='ID'
            value={updatedUser.ID}
            onChange={(e) =>
              setUpdateUser({ ...updatedUser, ID: e.target.value })
            }
            className='mb-2 w-full p-2 border border-gray-300 rounded'
          />
          <input
            type='text'
            placeholder='First Name'
            value={updatedUser.FirstName}
            onChange={(e) =>
              setUpdateUser({ ...updatedUser, FirstName: e.target.value })
            }
            className='mb-2 w-full p-2 border border-gray-300 rounded'
          />
          <input
            type='text'
            placeholder='Last Name'
            value={updatedUser.LastName}
            onChange={(e) =>
              setUpdateUser({ ...updatedUser, LastName: e.target.value })
            }
            className='mb-2 w-full p-2 border border-gray-300 rounded'
          />
          <input
            type='text'
            placeholder='Email'
            value={updatedUser.Email}
            onChange={(e) =>
              setUpdateUser({ ...updatedUser, Email: e.target.value })
            }
            className='mb-2 w-full p-2 border border-gray-300 rounded'
          />
          <button
            className='bg-blue-400 w-full p-2 rounded hover:bg-blue-500'
            type='submit'
          >
            Update User
          </button>
        </form>
      </div>

      <div className='flex flex-col gap-4 justify-center items-center space-y-4 bg-cyan-800 w-full rounded'>
        <h1 className='font-bold text-3xl py-4'>Users</h1>
        {users.map((user) => (
          <div key={user.ID} className='flex gap-2 w-full justify-between'>
            <CardComponent card={user} />
            <button
              onClick={() => deleteUser(user.ID)}
              className='bg-orange-400 rounded p-2'
            >
              Delete User
            </button>
          </div>
        ))}
      </div>
    </div>
  );
};

export default UserInterface;
