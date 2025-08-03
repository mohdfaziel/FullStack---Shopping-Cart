import { useState, useEffect } from 'react';
import { ordersAPI } from '../utils/api';

const OrderHistoryModal = ({ isOpen, onClose }) => {
  const [orders, setOrders] = useState([]);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    if (isOpen) {
      fetchOrders();
    }
  }, [isOpen]);

  const fetchOrders = async () => {
    setIsLoading(true);
    try {
      const data = await ordersAPI.getAll();
      setOrders(data || []);
    } catch (error) {
      console.error('Error fetching orders:', error);
      setOrders([]);
    } finally {
      setIsLoading(false);
    }
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

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleString();
  };

  const getStatusBadge = (status) => {
    const statusClasses = {
      'pending': 'badge-warning',
      'completed': 'badge-success',
      'cancelled': 'badge-error',
      'processing': 'badge-info'
    };
    
    return `badge ${statusClasses[status] || 'badge-neutral'} badge-sm`;
  };

  if (!isOpen) return null;

  return (
    <div className="modal modal-open">
      <div className="modal-box max-w-4xl">
        <div className="flex justify-between items-center mb-4">
          <h3 className="font-bold text-lg">📋 Order History</h3>
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
        ) : orders.length > 0 ? (
          <div>
            <div className="mb-4">
              <div className="stats shadow w-full">
                <div className="stat">
                  <div className="stat-title">Total Orders</div>
                  <div className="stat-value text-primary">{orders.length}</div>
                  <div className="stat-desc">All time orders placed</div>
                </div>
              </div>
            </div>

            <div className="space-y-4 max-h-96 overflow-y-auto">
              {orders.map((order) => (
                <div key={order.id} className="card bg-base-100 border shadow-sm">
                  <div className="card-body p-4">
                    <div className="flex justify-between items-start mb-3">
                      <div>
                        <h4 className="font-bold text-lg">Order #{order.id}</h4>
                        <p className="text-sm text-base-content/70">
                          Placed on {formatDate(order.created_at)}
                        </p>
                      </div>
                      <div className="text-right">
                        <span className={getStatusBadge(order.status)}>
                          {order.status || 'completed'}
                        </span>
                        {order.cart?.cart_items && order.cart.cart_items.length > 0 ? (
                          <p className="font-bold text-lg mt-1">
                            ₹{order.cart.cart_items.reduce((total, cartItem) => {
                              return total + (getItemPrice(cartItem.item.name) * cartItem.quantity);
                            }, 0).toLocaleString('en-IN')}
                          </p>
                        ) : order.total ? (
                          <p className="font-bold text-lg mt-1">₹{order.total.toLocaleString('en-IN')}</p>
                        ) : (
                          <p className="text-sm text-base-content/70">No pricing info</p>
                        )}
                      </div>
                    </div>

                    {order.cart?.cart_items && order.cart.cart_items.length > 0 ? (
                      <div>
                        <h5 className="font-semibold mb-2">Items Ordered:</h5>
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-2">
                          {order.cart.cart_items.map((cartItem, index) => (
                            <div key={index} className="bg-base-200 p-2 rounded flex justify-between items-center">
                              <div>
                                <span className="text-sm font-medium">{cartItem.item.name}</span>
                                <span className="text-xs text-base-content/70 ml-2">x{cartItem.quantity}</span>
                              </div>
                              <span className="text-sm font-medium">₹{(getItemPrice(cartItem.item.name) * cartItem.quantity).toLocaleString('en-IN')}</span>
                            </div>
                          ))}
                        </div>
                      </div>
                    ) : (
                      <div className="bg-base-200 p-3 rounded text-center">
                        <p className="text-sm text-base-content/70">
                          {order.cart ? 'No items in this order' : 'Order details not available'}
                        </p>
                      </div>
                    )}

                    <div className="flex justify-between items-center mt-3 pt-3 border-t">
                      <div className="text-sm text-base-content/70">
                        Cart ID: {order.cart_id || 'N/A'}
                      </div>
                      <div className="text-sm">
                        User ID: {order.user_id}
                      </div>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        ) : (
          <div className="text-center py-12">
            <div className="text-6xl mb-4">📋</div>
            <p className="text-lg font-medium mb-2">No orders found</p>
            <p className="text-base-content/70">You haven't placed any orders yet. Start shopping to see your order history here!</p>
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

export default OrderHistoryModal;
