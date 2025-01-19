'use client';

import Login from '@/components/Login';
import { useAuth } from '@/context/AuthContext';
import { useForm } from '@/context/FormContext';
import Register from '@/components/Register';
import AdminDashboard from '@/components/AdminDashboard';
import UserDashboard from '@/components/UserDashboard';

export default function Home() {
  const { isAuthenticated, isAdmin, isLoading } = useAuth();
  const { loginForm } = useForm();

  return (
    <div className='grid grid-rows-[20px_1fr_20px] items-center justify-items-center min-h-screen p-8 pb-20 gap-16 sm:p-20 font-[family-name:var(--font-geist-sans)] bg-[#D8CBAF]'>
      <main className='flex flex-col gap-8 row-start-2 items-center  w-full max-w-sm'>
        {isLoading ? (
          <div>Loading...</div>
        ) : (
          <>
            {isAuthenticated ? (
              isAdmin ? (
                <AdminDashboard />
              ) : (
                <UserDashboard />
              )
            ) : (
              <div className='w-full flex flex-col items-center'>
                <h1 className='text-[#4a4a4a] text-3xl font-extrabold'>
                  Welcome to The Nook
                </h1>
                <p className='text-[#4a4a4a] font-light text-center'>
                  Reserve games for you and your friends by signing up, or
                  logging in below
                </p>
                <div className='max-w-sm pt-4'>
                  {loginForm ? <Login /> : <Register />}
                </div>
              </div>
            )}
          </>
        )}
      </main>
    </div>
  );
}
