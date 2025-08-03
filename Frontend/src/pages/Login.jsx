import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import toast from 'react-hot-toast';
import { authAPI, auth } from '../utils/api';

const Login = ({ setIsAuthenticated }) => {
  const [formData, setFormData] = useState({
    username: '',
    password: ''
  });
  const [isLoading, setIsLoading] = useState(false);
  const navigate = useNavigate();

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsLoading(true);

    try {
      console.log('🔍 Login: Attempting to login with:', formData.username);
      const response = await authAPI.login(formData);
      console.log('🔍 Login: Login successful, response:', response);
      
      // Store token (this will automatically clear previous user's cart data)
      console.log('🔍 Login: Setting token and clearing previous cart data');
      auth.setToken(response.token);
      
      // Update authentication state
      setIsAuthenticated(true);
      
      // Show success toast
      toast.success('Login successful! Welcome back! 👋');
      
      // Navigate to items page
      navigate('/items');
    } catch (error) {
      console.error('🔍 Login: Login failed:', error);
      // Show error toast as specified in requirements
      toast.error('Invalid username/password');
    } finally {
      setIsLoading(false);
    }
  };

  const handleSignup = async () => {
    if (!formData.username || !formData.password) {
      toast.error('Please enter username and password');
      return;
    }

    setIsLoading(true);
    try {
      await authAPI.signup(formData);
      toast.success('User created successfully! Please login. ✅');
    } catch (error) {
      toast.error('Error creating user. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-base-200 flex items-center justify-center p-4">
      <div className="card w-full max-w-md bg-base-100 shadow-xl">
        <div className="card-body">
          <h2 className="card-title text-2xl font-bold text-center justify-center mb-6">
            Shopping Cart Login
          </h2>
          
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="form-control">
              <label className="label">
                <span className="label-text">Username</span>
              </label>
              <input
                type="text"
                name="username"
                value={formData.username}
                onChange={handleInputChange}
                placeholder="Enter your username"
                className="input input-bordered w-full"
                required
              />
            </div>

            <div className="form-control">
              <label className="label">
                <span className="label-text">Password</span>
              </label>
              <input
                type="password"
                name="password"
                value={formData.password}
                onChange={handleInputChange}
                placeholder="Enter your password"
                className="input input-bordered w-full"
                required
              />
            </div>

            <div className="form-control mt-6 space-y-2">
              <button 
                type="submit"
                className={`btn btn-primary w-full ${isLoading ? 'loading' : ''}`}
                disabled={isLoading}
              >
                {isLoading ? 'Logging in...' : 'Login'}
              </button>
              
              <button
                type="button"
                onClick={handleSignup}
                className={`btn btn-outline w-full ${isLoading ? 'loading' : ''}`}
                disabled={isLoading}
              >
                {isLoading ? 'Creating...' : 'Create Account'}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

export default Login;
