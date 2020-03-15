import pandas as pd
import statistics as stats

def make_dataset(path):

	dataset = pd.read_csv(path, header=None)
	X = dataset.iloc[:, 1:50].values
	y = dataset.iloc[:, 50].values
	return X, y

def multilayer_perceptron(X, y, name):
    from sklearn.neural_network import MLPClassifier
    clf = MLPClassifier(hidden_layer_sizes=(10, 10, 10), max_iter=1000000)
    kfold(clf, X, y, name)

def support_vector_machine(X, y, name):
    from sklearn.svm import SVC
    clf = SVC(kernel='linear')
    kfold(clf, X, y, name)

def k_nearest_neighbors(X, y, name):
    from sklearn.neighbors import KNeighborsClassifier
    clf = KNeighborsClassifier(n_neighbors=5)
    kfold(clf, X, y, name)

def random_forest(X, y, name):
    from sklearn.ensemble import RandomForestClassifier
    clf = RandomForestClassifier(n_estimators=20, random_state=0)
    kfold(clf, X, y, name)

def kfold(clf, X, y, name):
	from sklearn.preprocessing import StandardScaler
	from sklearn.model_selection import cross_val_predict

	if(name == "knn"):
		scaler = StandardScaler()
		scaler.fit(X)
		X = scaler.transform(X)

	y_pred = cross_val_predict(clf, X, y, cv=10)
	print_evaluation(y, y_pred)

def print_evaluation(y_test,y_pred):
    from sklearn.metrics import classification_report, confusion_matrix, accuracy_score
    print(confusion_matrix(y_test,y_pred))
    print(classification_report(y_test,y_pred))
    print("Average accuracy score:", accuracy_score(y_test, y_pred))

def main():
    import argparse
    parser = argparse.ArgumentParser()
    parser.add_argument("-m", "--mode", help="the ML mode to play with: mlp, svm, knn or rf are available")
    parser.add_argument("--path", help="path to the path csv")
    args = parser.parse_args()

    if args.path == None or args.mode == None:
        print("run with -h to know all the things")
        return

    X, y = make_dataset(args.path)

    if args.mode == "mlp":
        multilayer_perceptron(X, y, args.mode)
    elif args.mode == "svm":
        support_vector_machine(X, y, args.mode)
    elif args.mode == "knn":
        k_nearest_neighbors(X, y, args.mode)
    elif args.mode == "rf":
        random_forest(X, y, args.mode)
    else:
        print("run with -h to know all the things")

main()
