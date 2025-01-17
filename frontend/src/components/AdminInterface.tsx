import { useAuth } from '@/context/AuthContext';
import axios from 'axios';
import { useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';

const AdminInterface: React.FC = () => {
  const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
  const [users, setUsers] = useState<User[]>([]);
  const { isAuthenticated, name, jwt, email, logout } = useAuth();
  const router = useRouter();

  const handleLogout = async () => {
    await logout();
  };

  return (
    <div className='flex flex-col'>
      <p>HELLO ADMIN {name}</p>
      <button onClick={handleLogout}>Logout</button>
      <button onClick={() => router.push('/dashboard')}>Go to dashboard</button>
    </div>
  );
};

export default AdminInterface;
