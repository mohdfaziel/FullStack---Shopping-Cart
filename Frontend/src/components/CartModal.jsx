import { useState, useEffect } from 'react';
import toast from 'react-hot-toast';
import { cartAPI } from '../utils/api';

const CartModal = ({ isOpen, onClose }) => {
  const [cartData, setCartData] = useState(null);
  const [isLoading, setIsLoading] = useState(false);
  const [deletingItemId, setDeletingItemId] = useState(null);

  useEffect(() => {
    if (isOpen) {
      fetchCartData();
    }
  }, [isOpen]);

  // Add function to clear deleted items (for demo reset)
  const clearDeletedItems = () => {
    localStorage.removeItem('deletedCartItems');
    fetchCartData();
  };

  const fetchCartData = async () => {
    setIsLoading(true);
    try {
      const data = await cartAPI.get();
      setCartData(data);
    } catch (error) {
      console.error('Error fetching cart:', error);
      // If it's a 404 (no cart found), set empty cart structure
      if (error.message.includes('404')) {
        setCartData({ cart_items: [] });
      } else {
        // For other errors, set null to show error message
        setCartData(null);
      }
    } finally {
      setIsLoading(false);
    }
  };

  const processCartItems = () => {
    if (!cartData?.cart_items) return [];
    
    // Use cart_items directly since backend now properly manages cart state
    return cartData.cart_items.map(cartItem => ({
      ...cartItem.item,
      quantity: cartItem.quantity,
      cart_item_id: cartItem.id,
      item_id: cartItem.item_id
    }));
  };

  const calculateTotal = () => {
    const processedItems = processCartItems();
    return processedItems.reduce((total, item) => {
      const price = getItemPrice(item.name);
      return total + (price * item.quantity);
    }, 0);
  };

  const getItemPrice = (itemName) => {
    // Indian Rupee prices for demo purposes since backend doesn't have price field
    const prices = {
      'Laptop': 79999,
      'Smartphone': 49999,
      'Headphones': 7999,
      'Keyboard': 3999,
      'Mouse': 1999,
      'Monitor': 24999,
      'Tablet': 34999,
      'Webcam': 4999
    };
    return prices[itemName] || 9999;
  };

  const handleDeleteItem = async (itemId) => {
    setDeletingItemId(itemId);
    try {
      await cartAPI.removeItem(itemId);
      toast.success('Item removed from cart! 🗑️');
      
      // Refresh cart data (backend now properly manages state)
      await fetchCartData();
    } catch (error) {
      console.error('Error removing item from cart:', error);
      toast.error('Error removing item from cart');
    } finally {
      setDeletingItemId(null);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="modal modal-open">
      <div className="modal-box max-w-2xl">
        <div className="flex justify-between items-center mb-4">
          <h3 className="font-bold text-lg">🛒 Shopping Cart</h3>
          <button 
            className="btn btn-sm btn-circle btn-ghost"
            onClick={onClose}
          >
            ✕
          </button>
        </div>

        {isLoading ? (
          <div className="flex justify-center items-center py-8">
            <span className="loading loading-spinner loading-lg"></span>
          </div>
        ) : cartData ? (
          <div>
            <div className="mb-4">
              <div className="stat">
                <div className="stat-title">Cart ID</div>
                <div className="stat-value text-primary text-lg">{cartData.id}</div>
              </div>
            </div>

            {cartData.cart_items && cartData.cart_items.length > 0 ? (
              <div>
                <h4 className="font-semibold mb-3">Items in Cart:</h4>
                <div className="space-y-3 mb-6">
                  {processCartItems().map((item, index) => (
                    <div key={index} className="card bg-base-100 border">
                      <div className="card-body p-4">
                        <div className="flex justify-between items-start">
                          <div className="flex-1">
                            <h5 className="font-medium">{item.name}</h5>
                            <p className="text-sm text-base-content/70">
                              Status: <span className="badge badge-success badge-sm">{item.status}</span>
                            </p>
                            <p className="text-sm text-base-content/70">
                              Quantity: <span className="badge badge-info badge-sm">{item.quantity}</span>
                            </p>
                          </div>
                          <div className="flex items-center gap-3">
                            <div className="text-right">
                              <p className="font-bold text-lg">₹{(getItemPrice(item.name) * item.quantity).toLocaleString('en-IN')}</p>
                              <p className="text-xs text-base-content/60">
                                ₹{getItemPrice(item.name).toLocaleString('en-IN')} each
                              </p>
                              <p className="text-xs text-base-content/60">
                                Added: {new Date(item.created_at).toLocaleDateString()}
                              </p>
                            </div>
                            <button
                              onClick={() => handleDeleteItem(item.item_id)}
                              disabled={deletingItemId === item.item_id}
                              className={`btn btn-error btn-sm ${deletingItemId === item.item_id ? 'loading' : ''}`}
                              title="Remove from cart"
                            >
                              {deletingItemId === item.item_id ? 'Removing...' : '🗑️'}
                            </button>
                          </div>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>

                <div className="divider"></div>
                <div className="flex justify-between items-center">
                  <span className="text-lg font-semibold">Total:</span>
                  <span className="text-2xl font-bold text-primary">₹{calculateTotal().toLocaleString('en-IN')}</span>
                </div>
              </div>
            ) : (
              <div className="text-center py-8">
                <div className="text-6xl mb-4">🛒</div>
                <p className="text-lg font-medium mb-2">Your cart is empty</p>
                <p className="text-base-content/70">Add some items to get started!</p>
              </div>
            )}
          </div>
        ) : (
          <div className="text-center py-8">
            <div className="text-6xl mb-4">❌</div>
            <p className="text-lg font-medium mb-2">Unable to load cart</p>
            <p className="text-base-content/70">Please try again later</p>
          </div>
        )}

        <div className="modal-action">
          <button className="btn btn-primary" onClick={onClose}>
            Close
          </button>
        </div>
      </div>
    </div>
  );
};

export default CartModal;
