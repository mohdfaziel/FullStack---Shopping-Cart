import { Navigate } from 'react-router-dom';

const ProtectedRoute = ({ children, isAuthenticated, setIsAuthenticated }) => {
  return isAuthenticated ? children : <Navigate to="/login" replace />;
};

export default ProtectedRoute;
