import metapy


def main(config_path):
    invidx = metapy.index.make_inverted_index(config_path)


if __name__ == '__main__':
    main(config_path='config.toml')
