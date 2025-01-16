import { useAuth } from '@/context/AuthContext';
import axios from 'axios';
import { useEffect, useState } from 'react';

const AdminInterface: React.FC = () => {
  const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
  const [users, setUsers] = useState<User[]>([]);
  const { isAuthenticated, name, jwt, email } = useAuth();

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
  }, [apiUrl]);

  const deleteUser = async (ID: string) => {
    try {
      await axios.delete(`${apiUrl}/api/users/${ID}`);
      setUsers(users.filter((user) => user.ID !== ID));
    } catch (error) {
      console.error('Error deleting user:', error);
    }
  };

  return <h1>HELLO ADMIN {name}</h1>;
};

export default AdminInterface;
