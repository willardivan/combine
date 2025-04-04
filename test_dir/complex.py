#!/usr/bin/env python3
# This is a complex multiline Python file

import os
import sys
import json
from datetime import datetime

class ComplexExample:
    """
    This is a class with multiline docstring
    that should be compressed into a single line
    while maintaining readability.
    """
    
    def __init__(self, name, value):
        self.name = name
        self.value = value
        self.created_at = datetime.now()
        
    def process_data(self, data_list):
        """Process a list of data items."""
        result = []
        for item in data_list:
            if isinstance(item, dict):
                # Process dictionary items
                processed = {
                    "name": item.get("name", "unknown"),
                    "value": item.get("value", 0) * self.value,
                    "processed": True
                }
                result.append(processed)
            elif isinstance(item, list):
                # Process list items
                sub_result = []
                for sub_item in item:
                    sub_result.append(sub_item * 2)
                result.append(sub_result)
            else:
                # Default processing
                result.append(str(item) + "_" + self.name)
        
        return result

def main():
    # Create an instance of our class
    processor = ComplexExample("test", 10)
    
    # Sample data
    test_data = [
        {"name": "item1", "value": 5},
        ["a", "b", "c"],
        42,
        "hello"
    ]
    
    # Process the data
    result = processor.process_data(test_data)
    
    # Print results with nice formatting
    print(f"Processing results for {processor.name}:")
    for idx, item in enumerate(result):
        print(f"  [{idx}]: {item}")
    
    return 0

if __name__ == "__main__":
    sys.exit(main()) 