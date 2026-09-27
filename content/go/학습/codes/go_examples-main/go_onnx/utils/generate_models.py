import onnx
from onnx import helper
from onnx import TensorProto
import numpy as np
import os

# Create output directory
os.makedirs("models", exist_ok=True)

def create_simple_add_model(path):
    print(f"Generating {path}...")
    # Create input (ValueInfoProto)
    X = helper.make_tensor_value_info('X', TensorProto.FLOAT, [1])
    Y = helper.make_tensor_value_info('Y', TensorProto.FLOAT, [1])

    # Create output (ValueInfoProto)
    Z = helper.make_tensor_value_info('Z', TensorProto.FLOAT, [1])

    # Create a node (NodeProto)
    node_def = helper.make_node(
        'Add', # node name
        ['X', 'Y'], # inputs
        ['Z'], # outputs
    )

    # Create the graph (GraphProto)
    graph_def = helper.make_graph(
        [node_def],
        'test-model',
        [X, Y],
        [Z],
    )

    # Create the model (ModelProto)
    model_def = helper.make_model(graph_def, producer_name='onnx-example')
    
    onnx.checker.check_model(model_def)
    onnx.save(model_def, path)
    print(f"Saved {path}")

def create_iris_model(path):
    print(f"Generating {path}...")
    try:
        from sklearn.datasets import load_iris
        from sklearn.model_selection import train_test_split
        from sklearn.ensemble import RandomForestClassifier
        from skl2onnx import convert_sklearn
        from skl2onnx.common.data_types import FloatTensorType
        
        iris = load_iris()
        X, y = iris.data, iris.target
        X = X.astype(np.float32)
        X_train, X_test, y_train, y_test = train_test_split(X, y)
        
        clf = RandomForestClassifier(n_estimators=10)
        clf.fit(X_train, y_train)
        
        initial_type = [('float_input', FloatTensorType([None, 4]))]
        options = {'zipmap': False}
        onx = convert_sklearn(clf, initial_types=initial_type, options=options)
        
        with open(path, "wb") as f:
            f.write(onx.SerializeToString())
        print(f"Saved {path}")
    except ImportError as e:
        print(f"Skipping Iris model: {e} (install sklearn and skl2onnx)")

def create_mnist_model(path):
    print(f"Generating {path}...")
    try:
        import torch
        import torch.nn as nn
        import torch.nn.functional as F
        
        class Net(nn.Module):
            def __init__(self):
                super(Net, self).__init__()
                self.conv1 = nn.Conv2d(1, 10, kernel_size=5)
                self.conv2 = nn.Conv2d(10, 20, kernel_size=5)
                self.fc1 = nn.Linear(320, 50)
                self.fc2 = nn.Linear(50, 10)

            def forward(self, x):
                x = F.relu(F.max_pool2d(self.conv1(x), 2))
                x = F.relu(F.max_pool2d(self.conv2(x), 2))
                x = x.view(-1, 320)
                x = F.relu(self.fc1(x))
                x = self.fc2(x)
                return F.log_softmax(x, dim=1)

        model = Net()
        model.eval()
        
        # Dummy input for export (Batch size 1, 1 channel, 28x28 image)
        dummy_input = torch.randn(1, 1, 28, 28)
        
        torch.onnx.export(model, dummy_input, path, verbose=False,
                          input_names=['input'], output_names=['output'],
                          dynamic_axes={'input': {0: 'batch_size'}, 'output': {0: 'batch_size'}})
        print(f"Saved {path}")
        
    except ImportError as e:
        print(f"Skipping MNIST model: {e} (install torch)")


if __name__ == "__main__":
    create_simple_add_model("models/simple_add.onnx")
    create_iris_model("models/iris.onnx")
    create_mnist_model("models/mnist.onnx")
