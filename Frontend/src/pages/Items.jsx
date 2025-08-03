import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import toast from 'react-hot-toast';
import { itemsAPI, cartAPI, ordersAPI, auth } from '../utils/api';
import CartModal from '../components/CartModal';
import CheckoutModal from '../components/CheckoutModal';
import OrderHistoryModal from '../components/OrderHistoryModal';

const Items = ({ setIsAuthenticated }) => {
  const [items, setItems] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [addingItemId, setAddingItemId] = useState(null);
  const [isCartModalOpen, setIsCartModalOpen] = useState(false);
  const [isCheckoutModalOpen, setIsCheckoutModalOpen] = useState(false);
  const [isOrderHistoryModalOpen, setIsOrderHistoryModalOpen] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    fetchItems();
  }, []);

  const fetchItems = async () => {
    try {
      const response = await itemsAPI.getAll();
      setItems(response);
    } catch (error) {
      console.error('Error fetching items:', error);
      showErrorToast('Error loading items');
    } finally {
      setIsLoading(false);
    }
  };

  const handleAddToCart = async (itemId) => {
    setAddingItemId(itemId);
    try {
      // Add to backend first
      await cartAPI.addItem(itemId);
      
      toast.success('Item added to cart successfully! ✅');
      
      // Optionally update localStorage for immediate UI feedback
      // But prioritize backend as source of truth
      const currentCart = JSON.parse(localStorage.getItem('cart') || '[]');
      const itemExists = currentCart.find(item => item.item_id === itemId);
      
      if (itemExists) {
        // Increment quantity
        itemExists.quantity += 1;
      } else {
        // Add new item
        const itemData = items.find(item => item.id === itemId);
        if (itemData) {
          currentCart.push({
            id: Date.now(), // Temporary ID
            item_id: itemId,
            quantity: 1,
            item: itemData,
            created_at: new Date().toISOString()
          });
        }
      }
      
      // Save updated cart to localStorage
      localStorage.setItem('cart', JSON.stringify(currentCart));
      
    } catch (error) {
      console.error('Error adding item to cart:', error);
      toast.error('Failed to add item to cart. Please try again.');
    } finally {
      setAddingItemId(null);
    }
  };

  const handleCheckout = () => {
    setIsCheckoutModalOpen(true);
  };

  const handleViewCart = () => {
    setIsCartModalOpen(true);
  };

  const handleOrderHistory = () => {
    setIsOrderHistoryModalOpen(true);
  };

  // Get item price in Indian Rupees
  const getItemPrice = (itemName) => {
    const priceMap = {
      'Laptop': 79999,
      'Smartphone': 49999,
      'Headphones': 7999,
      'Keyboard': 3999,
      'Mouse': 1999,
      'Monitor': 24999,
      'Tablet': 34999,
      'Webcam': 4999
    };
    return priceMap[itemName] || 9999;
  };

  const handleCheckoutSuccess = () => {
    // Refresh items after successful checkout if needed
    fetchItems();
  };

  const handleLogout = () => {
    // Clear authentication and user data (this will clear cart automatically)
    auth.removeToken();
    setIsAuthenticated(false);
    navigate('/login');
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-base-200 flex items-center justify-center">
        <span className="loading loading-spinner loading-lg text-primary"></span>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-base-200">
      {/* Header */}
      <div className="navbar bg-base-100 shadow-md">
        <div className="flex-1">
          <h1 className="text-xl font-bold">Shopping Cart</h1>
        </div>
        <div className="flex-none">
          <button onClick={handleLogout} className="btn btn-ghost">
            Logout
          </button>
        </div>
      </div>

      <div className="container mx-auto p-6">
        {/* Action Buttons */}
        <div className="flex flex-wrap gap-4 mb-8 justify-center">
          <button 
            onClick={handleCheckout}
            className="btn btn-success"
          >
            💳 Checkout
          </button>
          <button 
            onClick={handleViewCart}
            className="btn btn-info"
          >
            🛒 Cart
          </button>
          <button 
            onClick={handleOrderHistory}
            className="btn btn-warning"
          >
            📋 Order History
          </button>
        </div>

        {/* Items Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
          {items.length === 0 ? (
            <div className="col-span-full text-center py-12">
              <h3 className="text-lg font-semibold text-base-content/70">
                No items available
              </h3>
            </div>
          ) : (
            items.map((item) => (
              <div key={item.id} className="card bg-base-100 shadow-xl hover:shadow-2xl transition-shadow">
                <div className="card-body">
                  <h2 className="card-title">
                    {item.name}
                    <div className={`badge ${item.status === 'available' ? 'badge-success' : 'badge-error'}`}>
                      {item.status}
                    </div>
                  </h2>
                  
                  <div className="text-lg font-bold text-primary mb-2">
                    ₹{getItemPrice(item.name).toLocaleString('en-IN')}
                  </div>
                  
                  <p className="text-base-content/70">
                    Created: {new Date(item.created_at).toLocaleDateString()}
                  </p>
                  
                  <div className="card-actions justify-end mt-4">
                    <button
                      onClick={() => handleAddToCart(item.id)}
                      disabled={item.status !== 'available' || addingItemId === item.id}
                      className={`btn btn-primary ${addingItemId === item.id ? 'loading' : ''}`}
                    >
                      {addingItemId === item.id ? 'Adding...' : 'Add to Cart'}
                    </button>
                  </div>
                </div>
              </div>
            ))
          )}
        </div>
      </div>

      {/* Modals */}
      <CartModal 
        isOpen={isCartModalOpen} 
        onClose={() => setIsCartModalOpen(false)} 
      />
      
      <CheckoutModal 
        isOpen={isCheckoutModalOpen} 
        onClose={() => setIsCheckoutModalOpen(false)}
        onSuccess={handleCheckoutSuccess}
      />
      
      <OrderHistoryModal 
        isOpen={isOrderHistoryModalOpen} 
        onClose={() => setIsOrderHistoryModalOpen(false)} 
      />
    </div>
  );
};

export default Items;
