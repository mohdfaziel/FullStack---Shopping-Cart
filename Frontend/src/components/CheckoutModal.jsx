import { useState, useEffect } from 'react';
import toast from 'react-hot-toast';
import { cartAPI, ordersAPI } from '../utils/api';

const CheckoutModal = ({ isOpen, onClose, onSuccess }) => {
  const [cartData, setCartData] = useState(null);
  const [isLoading, setIsLoading] = useState(false);
  const [isProcessing, setIsProcessing] = useState(false);

  useEffect(() => {
    if (isOpen) {
      fetchCartData();
    }
  }, [isOpen]);

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
    
    // Get deleted items from localStorage
    const deletedItems = JSON.parse(localStorage.getItem('deletedCartItems') || '[]');
    
    // Filter out deleted items and use cart_items directly since they now include quantities
    return cartData.cart_items
      .filter(cartItem => !deletedItems.includes(cartItem.item_id))
      .map(cartItem => ({
        ...cartItem.item,
        quantity: cartItem.quantity
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

  const handleCheckout = async () => {
    if (!cartData?.cart_items || cartData.cart_items.length === 0) {
      toast.error('Your cart is empty! Add items before checkout.');
      return;
    }

    setIsProcessing(true);
    try {
      await ordersAPI.create();
      toast.success('Order placed successfully! 🎉');
      // Clear deleted items after successful checkout
      localStorage.removeItem('deletedCartItems');
      onSuccess?.();
      onClose();
    } catch (error) {
      console.error('Error during checkout:', error);
      toast.error('Error processing order. Please try again.');
    } finally {
      setIsProcessing(false);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="modal modal-open">
      <div className="modal-box max-w-2xl">
        <div className="flex justify-between items-center mb-4">
          <h3 className="font-bold text-lg">💳 Checkout</h3>
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
        ) : cartData && cartData.cart_items && cartData.cart_items.length > 0 ? (
          <div>
            <div className="mb-6">
              <h4 className="font-semibold mb-3">Order Summary:</h4>
              <div className="space-y-2 mb-4">
                {processCartItems().map((item, index) => (
                  <div key={index} className="flex justify-between items-center p-3 bg-base-200 rounded">
                    <div>
                      <span className="font-medium">{item.name}</span>
                      <p className="text-sm text-base-content/70">{item.status}</p>
                      <p className="text-sm text-base-content/70">Qty: {item.quantity}</p>
                    </div>
                    <span className="font-bold">₹{(getItemPrice(item.name) * item.quantity).toLocaleString('en-IN')}</span>
                  </div>
                ))}
              </div>

              <div className="divider"></div>
              
              <div className="flex justify-between items-center mb-6">
                <span className="text-xl font-bold">Total Amount:</span>
                <span className="text-2xl font-bold text-primary">₹{calculateTotal().toLocaleString('en-IN')}</span>
              </div>

              <div className="bg-base-200 p-4 rounded-lg mb-4">
                <h5 className="font-semibold mb-2">📋 Order Details:</h5>
                <div className="grid grid-cols-2 gap-2 text-sm">
                  <div>Cart ID: <span className="font-medium">{cartData.id}</span></div>
                  <div>Items: <span className="font-medium">{processCartItems().length} unique ({cartData.cart_items.reduce((total, item) => total + item.quantity, 0)} total)</span></div>
                  <div>Status: <span className="badge badge-info badge-sm">Ready to Order</span></div>
                  <div>Date: <span className="font-medium">{new Date().toLocaleDateString()}</span></div>
                </div>
              </div>

              <div className="alert alert-info">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" className="stroke-current shrink-0 w-6 h-6">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                </svg>
                <span>By proceeding, you confirm your order and agree to our terms.</span>
              </div>
            </div>
          </div>
        ) : (
          <div className="text-center py-8">
            <div className="text-6xl mb-4">🛒</div>
            <p className="text-lg font-medium mb-2">Your cart is empty</p>
            <p className="text-base-content/70">Add items to your cart before checking out.</p>
          </div>
        )}

        <div className="modal-action">
          <button 
            className="btn btn-ghost" 
            onClick={onClose}
            disabled={isProcessing}
          >
            Cancel
          </button>
          {cartData && cartData.cart_items && cartData.cart_items.length > 0 && (
            <button 
              className="btn btn-success"
              onClick={handleCheckout}
              disabled={isProcessing}
            >
              {isProcessing ? (
                <>
                  <span className="loading loading-spinner loading-sm"></span>
                  Processing...
                </>
              ) : (
                <>
                  💳 Place Order (₹{calculateTotal().toLocaleString('en-IN')})
                </>
              )}
            </button>
          )}
        </div>
      </div>
    </div>
  );
};

export default CheckoutModal;
