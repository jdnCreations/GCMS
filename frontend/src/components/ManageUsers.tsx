import { useAuth } from '@/context/AuthContext';
import axios from 'axios';
import React, { useEffect, useState } from 'react';

const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

interface User {
  ID: string;
  FirstName: string;
  LastName: string;
  Email: string;
  IsAdmin: boolean;
}

interface ManageUsersProps {
  onCancel: () => void;
}

export default function ManageUsers({ onCancel }: ManageUsersProps) {
  const [users, setUsers] = useState<User[]>([]);
  const [showEditUser, setShowEditUser] = useState(false);
  const [selectedUser, setSelectedUser] = useState<User | null>(null);
  const [updatedUser, setUpdatedUser] = useState<User>({
    ID: '',
    FirstName: '',
    LastName: '',
    Email: '',
    IsAdmin: false,
  });
  const { jwt } = useAuth();

  const handleDeleteUser = async (userid: string) => {
    try {
      await axios.delete(`${apiUrl}/api/users/${userid}`, {
        headers: {
          Authorization: `Bearer ${jwt}`,
        },
      });
      setUsers(users.filter((user) => user.ID != userid));
    } catch (error) {
      console.log(error);
    }
  };

  const handleEditUser = (user: User) => {
    setSelectedUser(user);
    setUpdatedUser(user);
    setShowEditUser(true);
  };

  const handleUpdateUser = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!updatedUser.ID) {
      console.error('No user selected for update.');
      return;
    }
    try {
      await axios.put(`${apiUrl}/api/users/${updatedUser.ID}`, updatedUser, {
        headers: {
          Authorization: `Bearer ${jwt}`,
        },
      });
      setUsers((prevUsers) =>
        prevUsers.map((user) =>
          user.ID === updatedUser.ID ? updatedUser : user
        )
      );
      setShowEditUser(false);
    } catch (error) {
      console.error(error);
    }
  };

  useEffect(() => {
    const getUsers = async () => {
      try {
        const response = await axios.get(`${apiUrl}/api/users`, {
          headers: {
            Authorization: `Bearer ${jwt}`,
          },
        });
        setUsers(response.data);
      } catch (error) {
        console.log(error);
      }
    };
    getUsers();
  }, [jwt]);
  return (
    <div className='flex flex-col my-1 p-2 gap-1 bg-nook-charcoal rounded'>
      {users?.map((user) => (
        <div
          className='bg-nook-light-charcoal rounded px-2 py-1 flex justify-between'
          key={user.ID}
        >
          <p>
            {user.FirstName} {user.LastName}
          </p>
          <div className='flex gap-2'>
            <button
              onClick={() => handleEditUser(user)}
              className='hover:bg-nook-light-olive bg-nook-olive rounded px-1'
            >
              Update User
            </button>
            <button
              onClick={() => handleDeleteUser(user.ID)}
              className='hover:bg-nook-rose bg-nook-dark-rose rounded px-1'
            >
              Delete User
            </button>
          </div>
        </div>
      ))}
      {showEditUser && selectedUser && (
        <form onSubmit={(e) => handleUpdateUser(e)}>
          <p>
            Updating {selectedUser.FirstName} {selectedUser.LastName}
          </p>
          <input
            type='text'
            placeholder='First Name'
            value={updatedUser?.FirstName}
            onChange={(e) =>
              setUpdatedUser({
                ...updatedUser,
                FirstName: e.target.value,
                ID: selectedUser.ID,
              })
            }
          />
          <input
            type='text'
            placeholder='Last Name'
            value={updatedUser?.LastName}
            onChange={(e) =>
              setUpdatedUser({ ...updatedUser, LastName: e.target.value })
            }
          />
          <input
            type='text'
            placeholder='Email'
            value={updatedUser?.Email}
            onChange={(e) =>
              setUpdatedUser({ ...updatedUser, Email: e.target.value })
            }
          />
          <label htmlFor='isAdmin'>Admin</label>
          <input
            type='checkbox'
            id='isAdmin'
            checked={updatedUser?.IsAdmin || false}
            onChange={(e) => {
              setUpdatedUser({ ...updatedUser, IsAdmin: e.target.checked });
            }}
          />
          <button type='submit'>Update {selectedUser.FirstName}</button>
        </form>
      )}
      <button className='bg-nook-dark-rose rounded p-2' onClick={onCancel}>
        Cancel Managing Users
      </button>
    </div>
  );
}
